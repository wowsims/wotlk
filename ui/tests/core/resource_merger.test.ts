/// <reference types='jest' />
import { ResourceChangedLogGroup } from '../../core/proto_utils/logs_parser'
import { ResourceType } from '../../core/proto/api';
import { ResourceMerger } from '../../core/resources/merger';
import { emptySimLogParams } from '../fixtures/empty_sim_log_param';

const emptyResourceChangedLogGroupArray: ResourceChangedLogGroup[] = [];
const health_event_a = new ResourceChangedLogGroup(emptySimLogParams, ResourceType.ResourceTypeHealth, 100, 90, []);
const health_event_b = new ResourceChangedLogGroup(emptySimLogParams, ResourceType.ResourceTypeHealth, 100, 90, []);
const bloodRuneEventA = new ResourceChangedLogGroup(emptySimLogParams, ResourceType.ResourceTypeBloodRune, 2, 1, []);
const bloodRuneEventB = new ResourceChangedLogGroup(emptySimLogParams, ResourceType.ResourceTypeBloodRune, 1, 2, []);
const frostRuneEventA = new ResourceChangedLogGroup(emptySimLogParams, ResourceType.ResourceTypeFrostRune, 2, 1, []);
const frostRuneEventB = new ResourceChangedLogGroup(emptySimLogParams, ResourceType.ResourceTypeFrostRune, 1, 0, []);
const unholyRuneEventA = new ResourceChangedLogGroup(emptySimLogParams, ResourceType.ResourceTypeUnholyRune, 2, 1, []);
const unholyRuneEventB = new ResourceChangedLogGroup(emptySimLogParams, ResourceType.ResourceTypeUnholyRune, 1, 0, []);
const deathRuneEventA = new ResourceChangedLogGroup(emptySimLogParams, ResourceType.ResourceTypeUnholyRune, 0, 1, []);
const deathRuneEventB = new ResourceChangedLogGroup(emptySimLogParams, ResourceType.ResourceTypeUnholyRune, 1, 0, []);

describe('Resource merger', () => {
	it('Merges resources', () => {
		const input: Record<ResourceType, ResourceChangedLogGroup[]> = {
			[ResourceType.ResourceTypeNone]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeMana]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeEnergy]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeRage]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeComboPoints]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeFocus]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeRunicPower]: emptyResourceChangedLogGroupArray,
			[ResourceType.ResourceTypeHealth]: [health_event_a, health_event_b],
			[ResourceType.ResourceTypeBloodRune]: [bloodRuneEventA, bloodRuneEventB],
			[ResourceType.ResourceTypeFrostRune]: [frostRuneEventA, frostRuneEventB],
			[ResourceType.ResourceTypeUnholyRune]: [unholyRuneEventA, unholyRuneEventB],
			[ResourceType.ResourceTypeDeathRune]: [deathRuneEventA, deathRuneEventB],
		}
		const expected = input[ResourceType.ResourceTypeBloodRune]
			.concat(input[ResourceType.ResourceTypeFrostRune])
			.concat(input[ResourceType.ResourceTypeUnholyRune])
			.concat(input[ResourceType.ResourceTypeDeathRune]);
		const merger = new ResourceMerger([
			ResourceType.ResourceTypeBloodRune,
			ResourceType.ResourceTypeFrostRune,
			ResourceType.ResourceTypeUnholyRune,
			ResourceType.ResourceTypeDeathRune
		], ResourceType.ResourceTypeBloodRune);
		const output = merger.mergeResources(input);

		expect(output[ResourceType.ResourceTypeHealth]).toStrictEqual(input[ResourceType.ResourceTypeHealth]);
		expect(output[ResourceType.ResourceTypeBloodRune]).toStrictEqual(expected);
		expect(output[ResourceType.ResourceTypeFrostRune]).toStrictEqual([]);
		expect(output[ResourceType.ResourceTypeUnholyRune]).toStrictEqual([]);
		expect(output[ResourceType.ResourceTypeDeathRune]).toStrictEqual([]);
	})
})
