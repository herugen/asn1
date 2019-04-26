package def2

import (
	"reflect"

	"bitbucket.org/vservices/dark/logger"
	"github.com/jansemmelink/asn1/parser"
)

//IBlock ...
type IBlock interface {
	IParsable
	Start() IToken
	End() IToken
}

//Block is a parsable between a start and end keyword
type Block struct {
	Seq
}

//Start ...
func (Block) Start() IToken {
	return NewToken("") //will panic
}

//End ...
func (Block) End() IToken {
	return NewToken("") //will panic
}

//Parse ...
func (b Block) Parse(log logger.ILogger, l parser.ILines, v IParsable) (parser.ILines, error) {
	remain := l
	//structType := reflect.TypeOf(v).Elem()

	var userBlock IBlock
	userBlock = reflect.ValueOf(v).Elem().Interface().(IBlock)
	//log.Debugf("=====[ PARSING BLOCK %T: %s .. %s ]=====", b, b.Start(), b.End())   <<<< default block values not used
	log.Debugf("=====[ PARSING BLOCK %T: %s .. %s ]=====", userBlock, userBlock.Start(), userBlock.End())

	//skip over start token
	var err error
	if remain, err = userBlock.Start().Parse(log, remain, nil); err != nil {
		return l, log.Wrapf(nil, "%T%s...%s start token not found in line %d: %.32s ...", userBlock, userBlock.Start(), userBlock.End(), remain.LineNr(), remain.Next())
	}

	//parse contents the same way we parse any other struct as a Seq
	if remain, err = b.Seq.Parse(log, remain, v); err != nil {
		return l, log.Wrapf(err, "Failed to parse %T%s..%s contents from %d: %.32s ...", userBlock, userBlock.Start(), userBlock.End(), remain.LineNr(), remain.Next())
	}

	//skip over end token
	if remain, err = userBlock.End().Parse(log, remain, nil); err != nil {
		return l, log.Wrapf(nil, "%T%s...%s end token not found in line %d: %.32s ...", userBlock, userBlock.Start(), userBlock.End(), remain.LineNr(), remain.Next())
	}

	log.Debugf("Parsed %T%s...%s", userBlock, userBlock.Start(), userBlock.End())
	return remain, nil
} //Block.Parse()
