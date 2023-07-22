import { SimResult, SimResultFilter } from '../..//proto_utils/sim_result.js';
import { Component } from '../../components/component.js';
import { EventID, TypedEvent } from '../../typed_event.js';

export interface SimResultData {
	eventID: EventID,
	result: SimResult,
	filter: SimResultFilter,
};

export interface ResultComponentConfig {
	parent: HTMLElement,
	rootCssClass?: string,
	cssScheme?: String | null,
	resultsEmitter: TypedEvent<SimResultData | null>,
};

export abstract class ResultComponent extends Component {
	private lastSimResult: SimResultData | null;

	constructor(config: ResultComponentConfig) {
		super(config.parent, config.rootCssClass || 'result-component');
		this.lastSimResult = null;

		config.resultsEmitter.on((eventID, resultData) => {
			if (!resultData)
				return;

			this.lastSimResult = resultData;
			this.onSimResult(resultData);
		});
	}

	hasLastSimResult(): boolean {
		return this.lastSimResult != null;
	}

	getLastSimResult(): SimResultData {
		if (this.lastSimResult) {
			return this.lastSimResult;
		} else {
			throw new Error('No last sim result!');
		}
	}

	abstract onSimResult(resultData: SimResultData): void;
}
