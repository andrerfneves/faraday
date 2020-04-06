// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: rpc.proto

/*
Package frdrpc is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package frdrpc

import (
	"context"
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray

var (
	filter_FaradayServer_OutlierRecommendations_0 = &utilities.DoubleArray{Encoding: map[string]int{"rec_request": 0, "metric": 1}, Base: []int{1, 1, 1, 0}, Check: []int{0, 1, 2, 3}}
)

func request_FaradayServer_OutlierRecommendations_0(ctx context.Context, marshaler runtime.Marshaler, client FaradayServerClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq OutlierRecommendationsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		e   int32
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["rec_request.metric"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "rec_request.metric")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "rec_request.metric", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "rec_request.metric", err)
	}

	protoReq.RecRequest.Metric = CloseRecommendationRequest_Metric(e)

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_FaradayServer_OutlierRecommendations_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.OutlierRecommendations(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

var (
	filter_FaradayServer_ThresholdRecommendations_0 = &utilities.DoubleArray{Encoding: map[string]int{"rec_request": 0, "metric": 1}, Base: []int{1, 1, 1, 0}, Check: []int{0, 1, 2, 3}}
)

func request_FaradayServer_ThresholdRecommendations_0(ctx context.Context, marshaler runtime.Marshaler, client FaradayServerClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ThresholdRecommendationsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		e   int32
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["rec_request.metric"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "rec_request.metric")
	}

	err = runtime.PopulateFieldFromPath(&protoReq, "rec_request.metric", val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "rec_request.metric", err)
	}

	protoReq.RecRequest.Metric = CloseRecommendationRequest_Metric(e)

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_FaradayServer_ThresholdRecommendations_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.ThresholdRecommendations(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

var (
	filter_FaradayServer_RevenueReport_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_FaradayServer_RevenueReport_0(ctx context.Context, marshaler runtime.Marshaler, client FaradayServerClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq RevenueReportRequest
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_FaradayServer_RevenueReport_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.RevenueReport(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func request_FaradayServer_ChannelInsights_0(ctx context.Context, marshaler runtime.Marshaler, client FaradayServerClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ChannelInsightsRequest
	var metadata runtime.ServerMetadata

	msg, err := client.ChannelInsights(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// RegisterFaradayServerHandlerFromEndpoint is same as RegisterFaradayServerHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterFaradayServerHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterFaradayServerHandler(ctx, mux, conn)
}

// RegisterFaradayServerHandler registers the http handlers for service FaradayServer to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterFaradayServerHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterFaradayServerHandlerClient(ctx, mux, NewFaradayServerClient(conn))
}

// RegisterFaradayServerHandlerClient registers the http handlers for service FaradayServer
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "FaradayServerClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "FaradayServerClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "FaradayServerClient" to call the correct interceptors.
func RegisterFaradayServerHandlerClient(ctx context.Context, mux *runtime.ServeMux, client FaradayServerClient) error {

	mux.Handle("GET", pattern_FaradayServer_OutlierRecommendations_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_FaradayServer_OutlierRecommendations_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FaradayServer_OutlierRecommendations_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_FaradayServer_ThresholdRecommendations_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_FaradayServer_ThresholdRecommendations_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FaradayServer_ThresholdRecommendations_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_FaradayServer_RevenueReport_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_FaradayServer_RevenueReport_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FaradayServer_RevenueReport_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_FaradayServer_ChannelInsights_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_FaradayServer_ChannelInsights_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_FaradayServer_ChannelInsights_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_FaradayServer_OutlierRecommendations_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 1, 0, 4, 1, 5, 3}, []string{"v1", "faraday", "outliers", "rec_request.metric"}, ""))

	pattern_FaradayServer_ThresholdRecommendations_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 1, 0, 4, 1, 5, 3}, []string{"v1", "faraday", "threshold", "rec_request.metric"}, ""))

	pattern_FaradayServer_RevenueReport_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"v1", "faraday", "revenue"}, ""))

	pattern_FaradayServer_ChannelInsights_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"v1", "faraday", "insights"}, ""))
)

var (
	forward_FaradayServer_OutlierRecommendations_0 = runtime.ForwardResponseMessage

	forward_FaradayServer_ThresholdRecommendations_0 = runtime.ForwardResponseMessage

	forward_FaradayServer_RevenueReport_0 = runtime.ForwardResponseMessage

	forward_FaradayServer_ChannelInsights_0 = runtime.ForwardResponseMessage
)
