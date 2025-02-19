// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package dummy

import (
	"encoding/base64"
	"fmt"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_rpc "github.com/gagliardetto/solana-go/rpc"
	ag_base58 "github.com/mr-tron/base58"
	"reflect"
	"strings"
)

type MyEventEventData struct {
	Data  uint64
	Label string
}

var MyEventEventDataDiscriminator = [8]byte{96, 184, 197, 243, 139, 2, 90, 148}

func (obj MyEventEventData) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(MyEventEventDataDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `Data` param:
	err = encoder.Encode(obj.Data)
	if err != nil {
		return err
	}
	// Serialize `Label` param:
	err = encoder.Encode(obj.Label)
	if err != nil {
		return err
	}
	return nil
}

func (obj *MyEventEventData) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(MyEventEventDataDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %s, got %s",
				"[96 184 197 243 139 2 90 148]",
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `Data`:
	err = decoder.Decode(&obj.Data)
	if err != nil {
		return err
	}
	// Deserialize `Label`:
	err = decoder.Decode(&obj.Label)
	if err != nil {
		return err
	}
	return nil
}

func (*MyEventEventData) isEventData() {}

type MyOtherEventEventData struct {
	Data  uint64
	Label string
}

var MyOtherEventEventDataDiscriminator = [8]byte{133, 194, 247, 98, 118, 178, 28, 182}

func (obj MyOtherEventEventData) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(MyOtherEventEventDataDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `Data` param:
	err = encoder.Encode(obj.Data)
	if err != nil {
		return err
	}
	// Serialize `Label` param:
	err = encoder.Encode(obj.Label)
	if err != nil {
		return err
	}
	return nil
}

func (obj *MyOtherEventEventData) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(MyOtherEventEventDataDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %s, got %s",
				"[133 194 247 98 118 178 28 182]",
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `Data`:
	err = decoder.Decode(&obj.Data)
	if err != nil {
		return err
	}
	// Deserialize `Label`:
	err = decoder.Decode(&obj.Label)
	if err != nil {
		return err
	}
	return nil
}

func (*MyOtherEventEventData) isEventData() {}

var eventTypes = map[[8]byte]reflect.Type{
	MyEventEventDataDiscriminator:      reflect.TypeOf(MyEventEventData{}),
	MyOtherEventEventDataDiscriminator: reflect.TypeOf(MyOtherEventEventData{}),
}
var eventNames = map[[8]byte]string{
	MyEventEventDataDiscriminator:      "MyEvent",
	MyOtherEventEventDataDiscriminator: "MyOtherEvent",
}
var (
	_ *strings.Builder = nil
)
var (
	_ *base64.Encoding = nil
)
var (
	_ *ag_binary.Decoder = nil
)
var (
	_ *ag_rpc.GetTransactionResult = nil
)
var (
	_ *ag_base58.Alphabet = nil
)

type Event struct {
	Name string
	Data EventData
}

type EventData interface {
	UnmarshalWithDecoder(decoder *ag_binary.Decoder) error
	isEventData()
}

const eventLogPrefix = "Program data: "

func DecodeEvents(txData *ag_rpc.GetTransactionResult, targetProgramId ag_solanago.PublicKey, getAddressTables func(altAddresses []ag_solanago.PublicKey) (tables map[ag_solanago.PublicKey]ag_solanago.PublicKeySlice, err error)) (evts []*Event, err error) {
	var tx *ag_solanago.Transaction
	if tx, err = txData.Transaction.GetTransaction(); err != nil {
		return
	}

	altAddresses := make([]ag_solanago.PublicKey, len(tx.Message.AddressTableLookups))
	for i, alt := range tx.Message.AddressTableLookups {
		altAddresses[i] = alt.AccountKey
	}
	if len(altAddresses) > 0 {
		var tables map[ag_solanago.PublicKey]ag_solanago.PublicKeySlice
		if tables, err = getAddressTables(altAddresses); err != nil {
			return
		}
		tx.Message.SetAddressTables(tables)
		if err = tx.Message.ResolveLookups(); err != nil {
			return
		}
	}

	var base64Binaries [][]byte
	logMessageEventBinaries, err := decodeEventsFromLogMessage(txData.Meta.LogMessages)
	if err != nil {
		return
	}

	emitedCPIEventBinaries, err := decodeEventsFromEmitCPI(txData.Meta.InnerInstructions, tx.Message.AccountKeys, targetProgramId)
	if err != nil {
		return
	}

	base64Binaries = append(base64Binaries, logMessageEventBinaries...)
	base64Binaries = append(base64Binaries, emitedCPIEventBinaries...)
	evts, err = parseEvents(base64Binaries)
	return
}

func decodeEventsFromLogMessage(logMessages []string) (eventBinaries [][]byte, err error) {
	for _, log := range logMessages {
		if strings.HasPrefix(log, eventLogPrefix) {
			eventBase64 := log[len(eventLogPrefix):]

			var eventBinary []byte
			if eventBinary, err = base64.StdEncoding.DecodeString(eventBase64); err != nil {
				err = fmt.Errorf("failed to decode logMessage event: %s", eventBase64)
				return
			}
			eventBinaries = append(eventBinaries, eventBinary)
		}
	}
	return
}

func decodeEventsFromEmitCPI(InnerInstructions []ag_rpc.InnerInstruction, accountKeys ag_solanago.PublicKeySlice, targetProgramId ag_solanago.PublicKey) (eventBinaries [][]byte, err error) {
	for _, parsedIx := range InnerInstructions {
		for _, ix := range parsedIx.Instructions {
			if accountKeys[ix.ProgramIDIndex] != targetProgramId {
				continue
			}

			var ixData []byte
			if ixData, err = ag_base58.Decode(ix.Data.String()); err != nil {
				return
			}
			eventBase64 := base64.StdEncoding.EncodeToString(ixData[8:])
			var eventBinary []byte
			if eventBinary, err = base64.StdEncoding.DecodeString(eventBase64); err != nil {
				return
			}
			eventBinaries = append(eventBinaries, eventBinary)
		}
	}
	return
}

func parseEvents(base64Binaries [][]byte) (evts []*Event, err error) {
	decoder := ag_binary.NewDecoderWithEncoding(nil, ag_binary.EncodingBorsh)

	for _, eventBinary := range base64Binaries {
		eventDiscriminator := ag_binary.TypeID(eventBinary[:8])
		if eventType, ok := eventTypes[eventDiscriminator]; ok {
			eventData := reflect.New(eventType).Interface().(EventData)
			decoder.Reset(eventBinary)
			if err = eventData.UnmarshalWithDecoder(decoder); err != nil {
				err = fmt.Errorf("failed to unmarshal event %s: %w", eventType.String(), err)
				return
			}
			evts = append(evts, &Event{
				Name: eventNames[eventDiscriminator],
				Data: eventData,
			})
		}
	}
	return
}
