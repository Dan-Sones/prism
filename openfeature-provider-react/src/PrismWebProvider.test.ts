import { OpenFeature, ProviderStatus } from '@openfeature/web-sdk'
import { PrismWebProvider } from './PrismWebProvider.js'

describe('PrismWebProvider', () => {
    let fetchMock: jest.Mock
    const baseUrl = 'https://example.com'

    afterEach(() => {
        jest.resetAllMocks()
    })

    test('Should pass', () => {
        expect(true).toBe(true)
    })

    it('should initialize and have ready status', async () => {
        fetchMock = jest.fn().mockResolvedValue({
            ok: true,
            json: async () => ({
                'feature-flag-1': true,
                'feature-flag-2': 'blue',
                'feature-flag-3': 42,
            }),
        } as Response)

        OpenFeature.setContext({ targetingKey: 'test-key' })
        await OpenFeature.setProviderAndWait(new PrismWebProvider(fetchMock, baseUrl))

        expect(fetchMock).toHaveBeenCalledTimes(1)

        const client = OpenFeature.getClient()
        expect(client.providerStatus).toBe(ProviderStatus.READY)
    })

    describe('success tests', () => {
        beforeEach(async () => {
            fetchMock = jest.fn().mockResolvedValue({
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
            await OpenFeature.setProviderAndWait(new PrismWebProvider(fetchMock, baseUrl))

            expect(fetchMock).toHaveBeenCalledTimes(1)

            const client = OpenFeature.getClient()
            expect(client.providerStatus).toBe(ProviderStatus.READY)
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
})
