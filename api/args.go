package api

import "github.com/gin-gonic/gin"

//withChildrenArgs
func streamFieldsArgs(ctx *gin.Context) (flowUUID string, streamUUID string, producerUUID string, consumerUUID string, writerUUID string) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.FlowUUID = ctx.DefaultQuery(aType.FlowUUID, aDefault.FlowUUID)
	args.StreamUUID = ctx.DefaultQuery(aType.StreamUUID, aDefault.StreamUUID)
	args.ProducerUUID = ctx.DefaultQuery(aType.ProducerUUID, aDefault.ProducerUUID)
	args.ConsumerUUID = ctx.DefaultQuery(aType.ConsumerUUID, aDefault.ConsumerUUID)
	args.WriterUUID = ctx.DefaultQuery(aType.WriterUUID, aDefault.WriterUUID)
	return args.FlowUUID, args.StreamUUID, args.ProducerUUID, args.ConsumerUUID, args.WriterUUID
}

//withChildrenArgs
func withFieldsArgs(ctx *gin.Context) (field string, value string) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.Field = ctx.DefaultQuery(aType.Field, aDefault.Field)
	args.Value = ctx.DefaultQuery(aType.Value, aDefault.Value)
	return args.Field, args.Value
}

//withChildrenArgs
func queryFields(ctx *gin.Context) (order string, limit string) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault //ASC or DESC
	args.Order = ctx.DefaultQuery(aType.Order, aDefault.Order)
	args.Limit = ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints)
	return args.Order, args.Limit
}

//withChildrenArgs
func withChildrenArgs(ctx *gin.Context) (withChildren bool, withPoints bool) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithChildren = ctx.DefaultQuery(aType.WithChildren, aDefault.WithChildren)
	args.WithPoints = ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints)
	withChildren, _ = toBool(args.WithChildren) //?with_children=true&points=true
	withPoints, _ = toBool(args.WithPoints)
	return withChildren, withPoints
}

//withConsumerArgs
func withConsumerArgs(ctx *gin.Context) (askResponse bool, askRefresh bool, write bool, updateProducer bool) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.AskRefresh = ctx.DefaultQuery(aType.AskRefresh, aDefault.AskRefresh)
	args.AskResponse = ctx.DefaultQuery(aType.AskResponse, aDefault.AskResponse)
	args.Write = ctx.DefaultQuery(aType.Write, aDefault.Write)
	args.UpdateProducer = ctx.DefaultQuery(aType.UpdateProducer, aDefault.UpdateProducer)
	askRefresh, _ = toBool(args.AskRefresh)
	askResponse, _ = toBool(args.AskResponse)
	write, _ = toBool(args.Write)
	updateProducer, _ = toBool(args.UpdateProducer)
	return askRefresh, askResponse, write, updateProducer
}
