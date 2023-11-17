import { IndividualSimUI } from '../individual_sim_ui.js';
import { BaseModal } from './base_modal.js';
import { Player } from '../core/player.js';

export function addGemSummaryAction(simUI: IndividualSimUI<any>) {
	simUI.addAction('Gem Summary', 'gem-summary-action', () => {
		new GemSummaryMenu(simUI);
	});
}

class GemSummaryMenu extends BaseModal {
	private readonly simUI: IndividualSimUI<any>;
	private readonly tableBody: HTMLElement;
	private readonly player: Player<any>;

	constructor(simUI: IndividualSimUI<any>) {
		super(simUI.rootElem, 'gem-summary-menu', { scrollContents: true, size: 'md' });
		this.simUI = simUI;
		this.player = simUI.player;

		this.header?.insertAdjacentHTML('afterbegin', '<h5 class="modal-title">Currently Socketed Gems</h5>');
		this.body.innerHTML = `
			<div class="gem-summary-table-container modal-scroll-table">
				<table class="gem-summary-table" style="width: 100%">
					<thead>
						<tr>
							<th>Gem Type</th>
							<th style="text-align: right">
								<span>Quantity</span>
							</th>
						</tr>
					</thead>
					<tbody></tbody>
				</table>
			</div>
		`;
		
		this.tableBody = this.rootElem.querySelector('.gem-summary-table tbody') as HTMLElement;
		this.updateTable();
	}

	private updateTable() {
		this.tableBody.innerHTML = ``;
		const fullGemList = this.player.getGear().getAllGems(this.player.isBlacksmithing());
		const gemCounts = {};

		for (const gem of fullGemList) {
			gemCounts[gem.name] = gemCounts[gem.name] ? gemCounts[gem.name] + 1 : 1;
		}

		for (const gemName of Object.keys(gemCounts)) {
			const row = document.createElement('tr');
			row.innerHTML = `
				<td>${gemName}</td>
				<td style="text-align: right">${gemCounts[gemName].toFixed(0)}</td>
			`;
			this.tableBody.appendChild(row);
		}
	}

}
