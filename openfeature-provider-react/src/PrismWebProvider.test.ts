import { OpenFeature, ProviderStatus, StandardResolutionReasons } from '@openfeature/web-sdk'
import { PrismWebProvider } from './PrismWebProvider.js'

describe('PrismWebProvider', () => {
    const baseUrl = 'https://example.com'

    beforeEach(() => {
        global.fetch = jest.fn()
    })

    afterEach(() => {
        jest.resetAllMocks()
    })

    it('should initialize and have ready status', async () => {
        ;(global.fetch as jest.Mock).mockResolvedValue({
            ok: true,
            json: async () => ({
                'feature-flag-1': true,
                'feature-flag-2': 'blue',
                'feature-flag-3': 42,
            }),
        } as Response)

        OpenFeature.setContext({ targetingKey: 'test-key' })
        await OpenFeature.setProviderAndWait(new PrismWebProvider(baseUrl))

        expect(global.fetch).toHaveBeenCalledTimes(1)
        expect(global.fetch).toHaveBeenCalledWith(`${baseUrl}/api/assignments/test-key`, {})

        const client = OpenFeature.getClient()
        expect(client.providerStatus).toBe(ProviderStatus.READY)
    })

    describe('success tests', () => {
        beforeEach(async () => {
            ;(global.fetch as jest.Mock).mockResolvedValue({
                ok: true,
                json: async () => ({
                    'feature-flag-bool': true,
                    'feature-flag-string': 'blue',
                    'feature-flag-int': 42,
                    'feature-flag-float': 3.14,
                    'feature-flag-obj': { key: 'value' },
                }),
            } as Response)

            OpenFeature.setContext({ targetingKey: 'test-key' })
            await OpenFeature.setProviderAndWait(new PrismWebProvider(baseUrl))

            expect(global.fetch).toHaveBeenCalledTimes(1)

            const client = OpenFeature.getClient()
            expect(client.providerStatus).toBe(ProviderStatus.READY)
        })

        afterEach(() => {
            OpenFeature.clearProviders()
        })

        it('should evaluate boolean flag', async () => {
            const client = OpenFeature.getClient()
            const result = await client.getBooleanValue('feature-flag-bool', false)
            expect(result).toBe(true)
        })

        it('should evaluate string flag', async () => {
            const client = OpenFeature.getClient()
            const result = await client.getStringValue('feature-flag-string', 'red')
            expect(result).toBe('blue')
        })

        it('should evaluate integer flag', async () => {
            const client = OpenFeature.getClient()
            const result = await client.getNumberValue('feature-flag-int', 0)
            expect(result).toBe(42)
        })

        it('should evaluate float flag', async () => {
            const client = OpenFeature.getClient()
            const result = await client.getNumberValue('feature-flag-float', 0.0)
            expect(result).toBeCloseTo(3.14)
        })

        it('should evaluate object flag', async () => {
            const client = OpenFeature.getClient()
            const result = await client.getObjectValue('feature-flag-obj', {})
            expect(result).toEqual({ key: 'value' })
        })
    })

    describe('context change tests', () => {
        afterEach(() => {
            OpenFeature.clearProviders()
        })

        it('should update flags on context change', async () => {
            let activeKey = 'initial-key'

            ;(global.fetch as jest.Mock).mockImplementation(
                async () =>
                    ({
                        ok: true,
                        json: async () => ({
                            'feature-flag-bool': activeKey === 'initial-key' ? true : false,
                        }),
                    }) as Response
            )

            OpenFeature.setContext({ targetingKey: 'initial-key' })
            const provider = new PrismWebProvider(baseUrl)
            await OpenFeature.setProviderAndWait(provider)

            expect(global.fetch).toHaveBeenCalledTimes(1)

            const client = OpenFeature.getClient()

            let result = await client.getBooleanValue('feature-flag-bool', false)
            expect(result).toBe(true)

            activeKey = 'updated-key'
            await OpenFeature.setContext({ targetingKey: 'updated-key' })

            let contextChangeResult = await client.getBooleanValue('feature-flag-bool', false)

            expect(global.fetch).toHaveBeenCalledTimes(2)
            expect(contextChangeResult).toBe(false)
        })

        it("should not update flags if targetingKey hasn't changed", async () => {
            ;(global.fetch as jest.Mock).mockResolvedValue({
                ok: true,
                json: async () => ({
                    'feature-flag-bool': true,
                }),
            } as Response)

            OpenFeature.setContext({ targetingKey: 'same-key' })
            const provider = new PrismWebProvider(baseUrl)
            await OpenFeature.setProviderAndWait(provider)

            expect(global.fetch).toHaveBeenCalledTimes(1)

            const client = OpenFeature.getClient()

            let result = await client.getBooleanValue('feature-flag-bool', false)
            expect(result).toBe(true)

            await OpenFeature.setContext({ targetingKey: 'same-key' })

            let contextChangeResult = await client.getBooleanValue('feature-flag-bool', false)

            expect(global.fetch).toHaveBeenCalledTimes(1)
            expect(contextChangeResult).toBe(true)
        })
    })

    describe('error handling tests', () => {
        describe('http failure', () => {
            beforeEach(async () => {
                ;(global.fetch as jest.Mock).mockResolvedValue({
                    ok: false,
                    statusText: 'Internal Server Error',
                } as Response)

                OpenFeature.setContext({ targetingKey: 'test-key' })
                await expect(
                    OpenFeature.setProviderAndWait(new PrismWebProvider(baseUrl))
                ).rejects.toThrow(
                    'Failed to initialize PrismWebProvider: Failed to fetch flags: Internal Server Error'
                )
            })

            afterEach(() => {
                OpenFeature.clearProviders()
            })

            it('should handle fetch errors during initialization', () => {
                const client = OpenFeature.getClient()
                expect(client.providerStatus).toBe(ProviderStatus.ERROR)
            })

            it('should resolve default value when the flag fetch fails', async () => {
                const client = OpenFeature.getClient()
                expect(client.providerStatus).toBe(ProviderStatus.ERROR)

                const boolValue = await client.getBooleanValue('non-existent-flag', true)
                expect(boolValue).toBe(true)
            })
        })

        describe('missing targetingKey', () => {
            beforeEach(async () => {
                OpenFeature.setContext({})
                await expect(
                    OpenFeature.setProviderAndWait(new PrismWebProvider(baseUrl))
                ).rejects.toThrow(
                    'PrismWebProvider requires a targetingKey in the evaluation context for initialization.'
                )
            })

            afterEach(() => {
                OpenFeature.clearProviders()
            })

            it('should throw error if targetingKey is missing during initialization', () => {
                const client = OpenFeature.getClient()
                expect(client.providerStatus).toBe(ProviderStatus.ERROR)
            })

            it('should resolve default values when targetingKey is missing in context', async () => {
                const client = OpenFeature.getClient()

                const boolValue = await client.getBooleanValue('non-existent-flag', true)
                expect(boolValue).toBe(true)
            })
        })

        describe('non-existent flag', () => {
            beforeEach(async () => {
                ;(global.fetch as jest.Mock).mockResolvedValue({
                    ok: true,
                    json: async () => ({}),
                } as Response)

                OpenFeature.setContext({ targetingKey: 'test-key' })
                await OpenFeature.setProviderAndWait(new PrismWebProvider(baseUrl))
            })

            it('should resolve default value for non-existent flag', async () => {
                const client = OpenFeature.getClient()

                const stringValue = await client.getStringDetails('non-existent-flag', 'default')
                expect(stringValue.reason).toBe(StandardResolutionReasons.DEFAULT)
                expect(stringValue.value).toBe('default')
            })
        })
    })
})
