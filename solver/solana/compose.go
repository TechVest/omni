package solana

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/log"

	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	_ "embed"
)

// Start starts a genesis solana node and returns a client, a funded private key and a stop function or an error.
// The dir parameter is the location of the docker compose.
func Start(ctx context.Context, composeDir string) (*rpc.Client, solana.PrivateKey, func(), error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute) // Allow 1 minute for edge case of pulling images.
	defer cancel()

	if !composeDown(ctx, composeDir) {
		return nil, nil, nil, errors.New("failure to clean up previous solana instance")
	}

	// Ensure ports are available
	port, err := getAvailablePort()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "get available port")
	}

	if err := writeComposeFile(composeDir, port, "stable"); err != nil {
		return nil, nil, nil, errors.Wrap(err, "write compose file")
	}

	for _, program := range Programs {
		if err := copyProgram(composeDir, program); err != nil {
			return nil, nil, nil, errors.Wrap(err, "copy program")
		}
	}

	log.Info(ctx, "Starting solana")

	out, err := execCmd(ctx, composeDir, "docker", "compose", "up", "-d", "--remove-orphans")
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "docker compose up: "+out)
	}

	endpoint := "http://localhost:" + port

	cl := rpc.New(endpoint)

	stop := func() { //nolint:contextcheck // Fresh context required for stopping.
		// Fresh stop context since above context might be canceled.
		stopCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		composeDown(stopCtx, composeDir)
	}

	// Wait up to 10 secs for RPC to be available
	const retry = 20
	for i := 0; i < retry; i++ {
		if i == retry-1 {
			stop()
			return nil, nil, nil, errors.New("wait for RPC timed out")
		}

		select {
		case <-ctx.Done():
			return nil, nil, nil, errors.Wrap(ctx.Err(), "timeout")
		case <-time.After(time.Second):
		}

		_, err := cl.GetHealth(ctx)
		if err == nil {
			break
		}

		if i > retry/2 {
			log.Warn(ctx, "Solana: waiting for RPC to be available", err)
		}
	}

	privKey, err := solana.PrivateKeyFromSolanaKeygenFile(filepath.Join(composeDir, "id.json"))
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "get private key")
	}

	log.Info(ctx, "Solana: RPC is available", "addr", endpoint)

	return cl, privKey, stop, nil
}

// Deploy deploys a program to the localhost compose network using hte default keypair.
// It returns the confirmed transaction result or an error.
func Deploy(ctx context.Context, cl *rpc.Client, composeDir string, program Program) (*rpc.GetTransactionResult, error) {
	programDir := filepath.Join("/root/.config/solana/", program.Name)
	soFile := filepath.Join(programDir, program.SOFile())

	_, err := cl.GetAccountInfoWithOpts(ctx, program.MustPublicKey(), &rpc.GetAccountInfoOpts{
		Commitment: rpc.CommitmentConfirmed,
	})
	if !errors.Is(err, rpc.ErrNotFound) {
		return nil, errors.New("program already deployed", "account", program.MustPublicKey())
	}

	args := []string{
		"compose", "exec", "solana",
		"solana", "program", "deploy", soFile,
		"--verbose",
		"--commitment", string(rpc.CommitmentConfirmed),
		"--url", "localhost",
		"--output", "json",
	}

	var stdOut, stdErr bytes.Buffer
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Dir = composeDir
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, "exec deploy", "stderr", stdErr.String(), "stdout", stdOut.String())
	}

	var resp struct {
		ProgramID string `json:"programId"`
		Signature string `json:"signature"`
	}
	if err := json.Unmarshal(stdOut.Bytes(), &resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal deploy response", "stdout", stdOut.String())
	}

	if resp.ProgramID != program.MustPublicKey().String() {
		return nil, errors.New("unexpected program ID", "expected", program.MustPublicKey(), "got", resp.ProgramID)
	}

	txSig, err := solana.SignatureFromBase58(resp.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "parse signature")
	}

	tx, err := AwaitConfirmedTransaction(ctx, cl, txSig)
	if err != nil {
		return nil, errors.Wrap(err, "await confirmed transaction")
	}

	info, err := cl.GetAccountInfoWithOpts(ctx, program.MustPublicKey(), &rpc.GetAccountInfoOpts{
		Commitment: rpc.CommitmentConfirmed,
	})
	if err != nil {
		return nil, errors.Wrap(err, "get account info")
	} else if !info.Value.Executable {
		return nil, errors.New("program not executable", "account", program.MustPublicKey())
	}

	return tx, nil
}

func copyProgram(composeDir string, program Program) error {
	targetDir := filepath.Join(composeDir, program.Name)
	if err := os.RemoveAll(targetDir); err != nil {
		return errors.Wrap(err, "remove target dir", "path", targetDir)
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return errors.Wrap(err, "mkdir", "path", targetDir)
	}

	// Copy <programDir>/<program>.so and <programDir>/<program>-keypair.json to <composeDir>/<programDir>
	if err := os.WriteFile(
		filepath.Join(targetDir, program.SOFile()),
		program.SharedObject,
		0644,
	); err != nil {
		return errors.Wrap(err, "write so file")
	}

	if err := os.WriteFile(
		filepath.Join(targetDir, program.KeyPairFile()),
		program.KeyPairJSON,
		0644,
	); err != nil {
		return errors.Wrap(err, "write so file")
	}

	return nil
}

// composeDown runs docker-compose down in the provided directory.
func composeDown(ctx context.Context, dir string) bool {
	if _, err := os.Stat(dir + "/compose.yaml"); os.IsNotExist(err) {
		return true
	}

	out, err := execCmd(ctx, dir, "docker", "compose", "down")
	if err != nil {
		log.Error(ctx, "Error: docker compose down", err, "out", out)
		return false
	}

	log.Debug(ctx, "Solana: docker compose down: ok")

	return true
}

func execCmd(ctx context.Context, dir string, cmd string, args ...string) (string, error) {
	c := exec.CommandContext(ctx, cmd, args...)
	c.Dir = dir

	out, err := c.CombinedOutput()
	if err != nil {
		return string(out), errors.Wrap(err, "exec", "out", string(out))
	}

	return string(out), nil
}

//go:embed compose.yaml.tmpl
var composeTpl []byte

func writeComposeFile(dir string, port, version string) error {
	tpl, err := template.New("").Parse(string(composeTpl))
	if err != nil {
		return errors.Wrap(err, "parse compose template")
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, struct {
		Port    string
		Version string
	}{
		Port:    port,
		Version: version,
	})
	if err != nil {
		return errors.Wrap(err, "execute compose template")
	}

	err = os.WriteFile(dir+"/compose.yaml", buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, "write compose file")
	}

	return nil
}

func getAvailablePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", errors.Wrap(err, "resolve addr")
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", errors.Wrap(err, "listen")
	}
	defer l.Close()

	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return "", errors.Wrap(err, "split host port")
	}

	return port, nil
}
