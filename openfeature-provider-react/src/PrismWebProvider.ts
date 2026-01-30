import {
    ErrorCode,
    OpenFeatureEventEmitter,
    ProviderEvents,
    StandardResolutionReasons,
    type EvaluationContext,
    type Hook,
    type JsonValue,
    type Provider,
    type ProviderMetadata,
    type ResolutionDetails,
} from '@openfeature/web-sdk'
import type { Flags } from './types.js'

export class PrismWebProvider implements Provider {
    hooks?: Hook<Record<string, unknown>>[]

    readonly metadata: ProviderMetadata = {
        name: 'PrismWebProvider',
    }

    readonly events = new OpenFeatureEventEmitter()

    private baseUrl: string
    private fetchOptions: RequestInit

    private flagsCache: Flags | null = null

    constructor(baseUrl: string, fetchOptions: RequestInit = {}) {
        this.baseUrl = baseUrl
        this.fetchOptions = fetchOptions
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
        const response = await fetch(
            `${this.baseUrl}/assignments?userId=${userId}`,
            this.fetchOptions
        )
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
        if (this.flagsCache === null) {
            return {
                value: defaultValue,
                variant: 'default',
                reason: StandardResolutionReasons.UNKNOWN,
                errorCode: ErrorCode.PROVIDER_NOT_READY,
            }
        }

        if (this.flagsCache.hasOwnProperty(flagKey)) {
            const value = this.flagsCache[flagKey] as T
            return {
                value,
                variant: 'on',
                reason: StandardResolutionReasons.TARGETING_MATCH,
            }
        } else {
            return {
                value: defaultValue,
                variant: 'default',
                reason: StandardResolutionReasons.DEFAULT,
            }
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
