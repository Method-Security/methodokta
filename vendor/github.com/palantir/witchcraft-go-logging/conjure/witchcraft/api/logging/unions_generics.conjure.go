// This file was generated by Conjure and should not be manually edited.

//go:build go1.18

package logging

import (
	"context"
	"fmt"
)

type DiagnosticWithT[T any] Diagnostic

func (u *DiagnosticWithT[T]) Accept(ctx context.Context, v DiagnosticVisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "generic":
		if u.generic == nil {
			return result, fmt.Errorf("field \"generic\" is required")
		}
		return v.VisitGeneric(ctx, *u.generic)
	case "threadDump":
		if u.threadDump == nil {
			return result, fmt.Errorf("field \"threadDump\" is required")
		}
		return v.VisitThreadDump(ctx, *u.threadDump)
	}
}

type DiagnosticVisitorWithT[T any] interface {
	VisitGeneric(ctx context.Context, v GenericDiagnostic) (T, error)
	VisitThreadDump(ctx context.Context, v ThreadDumpV1) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}

type RequestLogWithT[T any] RequestLog

func (u *RequestLogWithT[T]) Accept(ctx context.Context, v RequestLogVisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "v1":
		if u.v1 == nil {
			return result, fmt.Errorf("field \"v1\" is required")
		}
		return v.VisitV1(ctx, *u.v1)
	case "v2":
		if u.v2 == nil {
			return result, fmt.Errorf("field \"v2\" is required")
		}
		return v.VisitV2(ctx, *u.v2)
	}
}

type RequestLogVisitorWithT[T any] interface {
	VisitV1(ctx context.Context, v RequestLogV1) (T, error)
	VisitV2(ctx context.Context, v RequestLogV2) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}

type UnionEventLogWithT[T any] UnionEventLog

func (u *UnionEventLogWithT[T]) Accept(ctx context.Context, v UnionEventLogVisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "eventLog":
		if u.eventLog == nil {
			return result, fmt.Errorf("field \"eventLog\" is required")
		}
		return v.VisitEventLog(ctx, *u.eventLog)
	case "eventLogV2":
		if u.eventLogV2 == nil {
			return result, fmt.Errorf("field \"eventLogV2\" is required")
		}
		return v.VisitEventLogV2(ctx, *u.eventLogV2)
	case "beaconLog":
		if u.beaconLog == nil {
			return result, fmt.Errorf("field \"beaconLog\" is required")
		}
		return v.VisitBeaconLog(ctx, *u.beaconLog)
	}
}

type UnionEventLogVisitorWithT[T any] interface {
	VisitEventLog(ctx context.Context, v EventLogV1) (T, error)
	VisitEventLogV2(ctx context.Context, v EventLogV2) (T, error)
	VisitBeaconLog(ctx context.Context, v BeaconLogV1) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}

type WrappedLogV1PayloadWithT[T any] WrappedLogV1Payload

func (u *WrappedLogV1PayloadWithT[T]) Accept(ctx context.Context, v WrappedLogV1PayloadVisitorWithT[T]) (T, error) {
	var result T
	switch u.typ {
	default:
		if u.typ == "" {
			return result, fmt.Errorf("invalid value in union type")
		}
		return v.VisitUnknown(ctx, u.typ)
	case "serviceLogV1":
		if u.serviceLogV1 == nil {
			return result, fmt.Errorf("field \"serviceLogV1\" is required")
		}
		return v.VisitServiceLogV1(ctx, *u.serviceLogV1)
	case "requestLogV2":
		if u.requestLogV2 == nil {
			return result, fmt.Errorf("field \"requestLogV2\" is required")
		}
		return v.VisitRequestLogV2(ctx, *u.requestLogV2)
	case "traceLogV1":
		if u.traceLogV1 == nil {
			return result, fmt.Errorf("field \"traceLogV1\" is required")
		}
		return v.VisitTraceLogV1(ctx, *u.traceLogV1)
	case "eventLogV2":
		if u.eventLogV2 == nil {
			return result, fmt.Errorf("field \"eventLogV2\" is required")
		}
		return v.VisitEventLogV2(ctx, *u.eventLogV2)
	case "metricLogV1":
		if u.metricLogV1 == nil {
			return result, fmt.Errorf("field \"metricLogV1\" is required")
		}
		return v.VisitMetricLogV1(ctx, *u.metricLogV1)
	case "auditLogV2":
		if u.auditLogV2 == nil {
			return result, fmt.Errorf("field \"auditLogV2\" is required")
		}
		return v.VisitAuditLogV2(ctx, *u.auditLogV2)
	case "diagnosticLogV1":
		if u.diagnosticLogV1 == nil {
			return result, fmt.Errorf("field \"diagnosticLogV1\" is required")
		}
		return v.VisitDiagnosticLogV1(ctx, *u.diagnosticLogV1)
	}
}

type WrappedLogV1PayloadVisitorWithT[T any] interface {
	VisitServiceLogV1(ctx context.Context, v ServiceLogV1) (T, error)
	VisitRequestLogV2(ctx context.Context, v RequestLogV2) (T, error)
	VisitTraceLogV1(ctx context.Context, v TraceLogV1) (T, error)
	VisitEventLogV2(ctx context.Context, v EventLogV2) (T, error)
	VisitMetricLogV1(ctx context.Context, v MetricLogV1) (T, error)
	VisitAuditLogV2(ctx context.Context, v AuditLogV2) (T, error)
	VisitDiagnosticLogV1(ctx context.Context, v DiagnosticLogV1) (T, error)
	VisitUnknown(ctx context.Context, typ string) (T, error)
}
