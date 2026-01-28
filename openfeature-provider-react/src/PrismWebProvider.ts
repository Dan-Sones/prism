import {
    ClientProviderStatus,
    type AnyProviderEvent,
    type EvaluationContext,
    type Hook,
    type JsonValue,
    type Logger,
    type Paradigm,
    type Provider,
    type ProviderEventEmitter,
    type ProviderMetadata,
    type ResolutionDetails,
    type TrackingEventDetails,
} from '@openfeature/web-sdk'

export class PrismWebProvider implements Provider {
    hooks?: Hook<Record<string, unknown>>[]
    onContextChange?(
        oldContext: EvaluationContext,
        newContext: EvaluationContext
    ): Promise<void> | void {
        throw new Error('Method not implemented.')
    }
    resolveBooleanEvaluation(
        flagKey: string,
        defaultValue: boolean,
        context: EvaluationContext,
        logger: Logger
    ): ResolutionDetails<boolean> {
        throw new Error('Method not implemented.')
    }
    resolveStringEvaluation(
        flagKey: string,
        defaultValue: string,
        context: EvaluationContext,
        logger: Logger
    ): ResolutionDetails<string> {
        throw new Error('Method not implemented.')
    }
    resolveNumberEvaluation(
        flagKey: string,
        defaultValue: number,
        context: EvaluationContext,
        logger: Logger
    ): ResolutionDetails<number> {
        throw new Error('Method not implemented.')
    }
    resolveObjectEvaluation<T extends JsonValue>(
        flagKey: string,
        defaultValue: T,
        context: EvaluationContext,
        logger: Logger
    ): ResolutionDetails<T> {
        throw new Error('Method not implemented.')
    }
    metadata: ProviderMetadata
    runsOn?: Paradigm
    status?: ClientProviderStatus
    events?: ProviderEventEmitter<AnyProviderEvent, Record<string, unknown>>
    onClose?(): Promise<void> {
        throw new Error('Method not implemented.')
    }
    initialize?(context?: EvaluationContext): Promise<void> {
        throw new Error('Method not implemented.')
    }
    track?(
        trackingEventName: string,
        context: EvaluationContext,
        trackingEventDetails: TrackingEventDetails
    ): void {
        throw new Error('Method not implemented.')
    }
}
