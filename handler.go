package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func errorResponse(err error) gin.H {
	// The Error method in on type  error itself, its the only method on type error, it returns the error in string format
	return gin.H{"error": err.Error()}
}

// Input Log Data
type logData struct {
	Level      string `json:"level" binding:"required"`
	Message    string `json:"message" binding:"required"`
	ResourceId string `json:"resourceId" binding:"required"`
	Timestamp  string `json:"timestamp" binding:"required"`
	TraceId    string `json:"traceId" binding:"required"`
	SpanId     string `json:"spanId" binding:"required"`
	Commit     string `json:"commit" binding:"required"`
	Metadata   struct {
		ParentResourceId string `json:"parentResourceId" binding:"required"`
	} `json:"metadata" binding:"required"`
}

// Log data structure to be stored in DB
type Log struct {
	Level      string    `bson:"level"`
	Message    string    `bson:"message"`
	ResourceId string    `bson:"resourceId"`
	Timestamp  time.Time `bson:"timestamp"`
	TraceId    string    `bson:"traceId"`
	SpanId     string    `bson:"spanId"`
	Commit     string    `bson:"commit"`
	Metadata   struct {
		ParentResourceId string `bson:"parentResourceId"`
	} `bson:"metadata"`
}

// Uploading log data to DB
func (server *Server) logUpload(ctx *gin.Context) {
	// Validation of Input Log data
	var req logData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	doc := Log{
		Level:      req.Level,
		Message:    req.Message,
		ResourceId: req.ResourceId,
		TraceId:    req.TraceId,
		SpanId:     req.SpanId,
		Commit:     req.Commit,
	}
	doc.Metadata.ParentResourceId = req.Metadata.ParentResourceId

	// Conv time string to time value
	t, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	doc.Timestamp = t

	coll := server.client.Database("db").Collection("logs")
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// Rendering index html for homepage, web ui for searching logs
func (server *Server) homepage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Log-App",
	})
}

// Search Operation
func (server *Server) search(ctx *gin.Context) {
	coll := server.client.Database("db").Collection("logs")
	textSearch := ctx.PostForm("text-search")

	// Aggregation Pipeline Stages
	opts := options.Aggregate().SetMaxTime(3 * time.Second)
	searchStage := bson.D{{"$search", bson.D{{"index", "default"}, {"text", bson.D{{"query", textSearch}, {"path", bson.D{{"wildcard", "*"}}}}}}}}
	limitStage := bson.D{{"$limit", 5}}

	formContent := Log{
		Level:      ctx.PostForm("level"),
		Message:    ctx.PostForm("message"),
		ResourceId: ctx.PostForm("resourceId"),
		TraceId:    ctx.PostForm("traceId"),
		SpanId:     ctx.PostForm("spanId"),
		Commit:     ctx.PostForm("commit"),
	}
	formContent.Metadata.ParentResourceId = ctx.PostForm("metadata.parentResourceId")

	filter := bson.D{}

	// FILTER QUERIES
	// Not adding the field values in filter which are empty
	// Together they perform AND operation, that is, even if one doesn't satisfy it will return null
	// So it will combine multiple filters (BONUS)
	// Also implemented Date Range (BONUS)
	if formContent.Level != "" {
		filter = append(filter, bson.E{"level", formContent.Level})
	}
	if formContent.Message != "" && textSearch == "" {
		filter = append(filter, bson.E{"$text", bson.D{{"$search", formContent.Message}}})
	} else if formContent.Message != "" {
		filter = append(filter, bson.E{"message", formContent.Message})
	}
	if formContent.ResourceId != "" {
		filter = append(filter, bson.E{"resourceId", formContent.ResourceId})
	}
	if formContent.TraceId != "" {
		filter = append(filter, bson.E{"traceId", formContent.TraceId})
	}
	if formContent.SpanId != "" {
		filter = append(filter, bson.E{"spanId", formContent.SpanId})
	}
	if formContent.Commit != "" {
		filter = append(filter, bson.E{"commit", formContent.Commit})
	}
	if formContent.Metadata.ParentResourceId != "" {
		filter = append(filter, bson.E{"metadata.parentResourceId", formContent.Metadata.ParentResourceId})
	}
	if ctx.PostForm("timestamp-lower") != "" {
		// Conv time string to time value, so that we can compare
		timestampL, err := time.Parse(time.RFC3339, ctx.PostForm("timestamp-lower"))
		if err != nil {
			fmt.Println(err)
		}
		filter = append(filter, bson.E{"timestamp", bson.D{{"$lte", timestampL}}})
	}
	if ctx.PostForm("timestamp-higher") != "" {
		timestampG, err := time.Parse(time.RFC3339, ctx.PostForm("timestamp-higher"))
		if err != nil {
			fmt.Println(err)
		}
		filter = append(filter, bson.E{"timestamp", bson.D{{"$gte", timestampG}}})
	}
	fmt.Println(filter)

	var results []Log

	if textSearch != "" {
		cursor, err := coll.Aggregate(ctx, mongo.Pipeline{searchStage, bson.D{{"$match", filter}}, limitStage}, opts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		if err = cursor.All(context.TODO(), &results); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		if err = cursor.All(context.TODO(), &results); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	data, error := json.MarshalIndent(results, "", "\t")
	if error != nil {
		fmt.Println("JSON parse error: ", error)
		return
	}
	ctx.String(http.StatusOK, string(data))
	// ctx.HTML(http.StatusOK, "result.html", gin.H{
	// 	"data": results,
	// })
}
