// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element, fragment } from 'tsx-vanilla';

import { SimLog } from '../../proto_utils/logs_parser.js';
import { TypedEvent } from '../../typed_event.js';
import { BooleanPicker } from '../boolean_picker.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

export class LogRunner extends ResultComponent {
	private logsContainer: HTMLElement;

	private showDebug = false;

	readonly showDebugChangeEmitter = new TypedEvent<void>('Show Debug');

	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'log-runner-root';
		super(config)

		this.rootElem.appendChild(
			<>
				<div className="show-debug-container"></div>
				<table className="metrics-table log-runner-table">
					<thead>
						<tr className="metrics-table-header-row">
							<th>Time</th>
							<th>
								<div className="d-flex align-items-end">Event</div>
							</th>
						</tr>
					</thead>
					<tbody className="log-runner-logs"></tbody>
				</table>
			</>
		)
		this.logsContainer = this.rootElem.querySelector('.log-runner-logs')!;

		new BooleanPicker<LogRunner>(this.rootElem.querySelector('.show-debug-container')!, this, {
			extraCssClasses: ['show-debug-picker'],
			label: 'Show Debug Statements',
			inline: true,
			reverse: true,
			changedEvent: () => this.showDebugChangeEmitter,
			getValue: () => this.showDebug,
			setValue: (eventID, _logRunner, newValue) => {
				this.showDebug = newValue;
				this.showDebugChangeEmitter.emit(eventID);
			}
		});

		this.showDebugChangeEmitter.on(() => this.onSimResult(this.getLastSimResult()));
	}

	onSimResult(resultData: SimResultData): void {
		const logs = resultData.result.logs
		this.logsContainer.innerHTML = '';
		logs.
			filter(log => !log.isCastCompleted()).
			forEach(log => {
				const lineElem = document.createElement('span');
				lineElem.textContent = log.toString();
				if (log.raw.length > 0 && (this.showDebug || !log.raw.match(/.*\[DEBUG\].*/))) {
					this.logsContainer.appendChild(
						<tr>
							<td className="log-timestamp">{log.formattedTimestamp()}</td>
							<td className="log-event">{this.newEventFrom(log)}</td>
						</tr>
					);
				}
			});
	}

	private newEventFrom(log: SimLog): Element {
		const eventString = log.toString(false).trim();
		const wrapper = <span></span>;
		wrapper.innerHTML = eventString;
		return wrapper;
	}
}
