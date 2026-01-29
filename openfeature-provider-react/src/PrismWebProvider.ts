import {
    ClientProviderStatus,
    OpenFeatureEventEmitter,
    ProviderEvents,
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
import type { Flags } from './types.js'

export class PrismWebProvider implements Provider {
    hooks?: Hook<Record<string, unknown>>[]

    readonly metadata: ProviderMetadata = {
        name: 'PrismWebProvider',
    }

    readonly events = new OpenFeatureEventEmitter()

    private fetchInstance: typeof fetch
    private baseUrl: string

    private flagsCache: Flags | null = null

    constructor(fetchInstance: typeof fetch, baseUrl: string) {
        this.fetchInstance = fetchInstance
        this.baseUrl = baseUrl
    }

    async initialize?(context?: EvaluationContext): Promise<void> {
        if (context?.targetingKey === undefined) {
            throw new Error(
                'PrismWebProvider requires a targetingKey in the evaluation context for initialization.'
            )
        }

        try {
            await this.updateFlagsCache(context.targetingKey)
        } catch (error) {
            this.events.emit(ProviderEvents.Error)
            throw new Error(`Failed to initialize PrismWebProvider: ${(error as Error).message}`)
        }

        this.events.emit(ProviderEvents.Ready)
    }

    async updateFlagsCache(userId: string): Promise<void> {
        const flags = await this.fetchFlags(userId)
        this.flagsCache = flags
    }

    async fetchFlags(userId: string): Promise<Flags> {
        const response = await this.fetchInstance(`${this.baseUrl}/assignments?userId=${userId}`)
        if (!response.ok) {
            throw new Error(`Failed to fetch flags: ${response.statusText}`)
        }
        const flags = (await response.json()) as Flags
        return flags
    }

    onContextChange?(
        oldContext: EvaluationContext,
        newContext: EvaluationContext
    ): Promise<void> | void {
        if (oldContext.targetingKey !== newContext.targetingKey) {
            return this.updateFlagsCache(newContext.targetingKey as string)
        }
    }

    private evaluateFlag<T extends any>(flagKey: string, defaultValue: T): ResolutionDetails<T> {
        // TODO: Handle errors better
        // how do we get the user to just see the control but not log their events?

        if (this.flagsCache === null) {
            throw new Error('Flags cache is not initialized.')
        }

        if (this.flagsCache.hasOwnProperty(flagKey)) {
            const value = this.flagsCache[flagKey] as T
            return {
                value,
                variant: 'on',
                reason: 'TARGETING_MATCH',
            }
        } else {
            throw new Error(`Flag with key ${flagKey} not found.`)
        }
    }

    resolveBooleanEvaluation(flagKey: string, defaultValue: boolean): ResolutionDetails<boolean> {
        return this.evaluateFlag(flagKey, defaultValue)
    }

    resolveStringEvaluation(flagKey: string, defaultValue: string): ResolutionDetails<string> {
        return this.evaluateFlag(flagKey, defaultValue)
    }

    resolveNumberEvaluation(flagKey: string, defaultValue: number): ResolutionDetails<number> {
        return this.evaluateFlag(flagKey, defaultValue)
    }

    resolveObjectEvaluation<T extends JsonValue>(
        flagKey: string,
        defaultValue: T
    ): ResolutionDetails<T> {
        return this.evaluateFlag(flagKey, defaultValue)
    }
}
