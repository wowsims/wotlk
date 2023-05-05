import {ProgressMetrics} from '../proto/api.js';
import {GemColor, ItemSlot} from '../proto/common.js';
import {IndividualSimUI} from '../individual_sim_ui.js';
import {EventID, TypedEvent} from '../typed_event.js';

import {BaseModal} from './base_modal.js';
import {ResultsViewer} from './results_viewer.js';
import {EquippedItem} from "../proto_utils/equipped_item";
import {UIGem, UIItem} from "../proto/ui";
import {getEmptyGemSocketIconUrl} from "../proto_utils/gems";
import {ItemData} from "./gear_picker";
import {ActionId} from "../proto_utils/action_id";
import {BooleanPicker} from "./boolean_picker";
import {GemSelectorModal} from "./individual_sim_ui/bulk_tab"

export function addDroptimizerAction(simUI: IndividualSimUI<any>) {
    simUI.addAction('Droptimizer', 'ep-weights-action', () => {
        new DroptimizerMenu(simUI);
    });
}

class DroptimizerMenu extends BaseModal {
    private readonly simUI: IndividualSimUI<any>;
    private readonly container: HTMLElement;
    private readonly table: HTMLElement;
    private readonly tableBody: HTMLElement;
    private readonly resultsViewer: ResultsViewer;

    private isCanceled: boolean = false;
    private currentSim: number = 0;
    private totalSims: number = 0;
    private statsType: string;
    private autoGem: boolean = false;
    private defaultGems: UIGem[];
    readonly itemsChangedEmitter = new TypedEvent<void>();
    private containerGemPicker: HTMLElement;



    constructor(simUI: IndividualSimUI<any>) {
        super(simUI.rootElem, 'ep-weights-menu', {footer: true, scrollContents: true});
        this.simUI = simUI;
        this.statsType = 'U25';
        this.defaultGems = [UIGem.create(), UIGem.create(), UIGem.create(), UIGem.create()];

        this.header?.insertAdjacentHTML('afterbegin', '<h5 class="modal-title">Generate Highest Value Items</h5>');
        this.body.innerHTML = `
			<p>Select the Dungeon/Raid you want to sim</p>
			<div class="ep-weights-options row">
				<div class="col col-sm-3">
					<select class="ep-type-select form-select">
						<option value="U10">Ulduar 10</option>
						<option value="U25">Ulduar 25</option>
					</select>
				</div>
				<div class="show-all-stats-container col col-sm-3"></div>
			</div>
			<div class="results-ep-table-container modal-scroll-table">
				<div class="results-pending-overlay"></div>
				<table class="results-ep-table">
					<thead>
						<tr>
							<th>Item</th>
							<th class="damage-metrics type-weight">
								<span>DPS Diff</span>
							</th>
						</tr>
					</thead>
					<tbody>
					 <td>Example Item</td>
					 <td class="stdev-cell type-weight">
					 <span class="results-avg">322</span>
				     </td>
					 <td class="current-ep"></td>
					</tbody>
				</table>
			</div>
			<p></p>
			<p></p>
			<p></p>
			<div class="content-block"></div>
		`;
        this.footer!.innerHTML = `
			<button class="btn btn-secondary">
				<i class="fas"></i>
				Cancel Operation
			</button>
			<button class="btn btn-primary calc-weights">
				<i class="fas fa-calculator"></i>
				Calculate
			</button>
		`;
        const resultsElem = this.rootElem.querySelector('.results-pending-overlay') as HTMLElement;
        this.resultsViewer = new ResultsViewer(resultsElem);
        this.container = this.rootElem.querySelector('.results-ep-table-container') as HTMLElement;
        this.table = this.rootElem.querySelector('.results-ep-table') as HTMLElement;
        this.tableBody = this.rootElem.querySelector('.results-ep-table tbody') as HTMLElement;
        this.containerGemPicker = this.rootElem.querySelector('.content-block') as HTMLElement;


        const selectElem = this.rootElem.getElementsByClassName('ep-type-select')[0] as HTMLSelectElement;
        selectElem.addEventListener('input', () => {
            this.statsType = selectElem.value;
        });
        selectElem.value = this.statsType;
        let resultList = Array<Result>();
        const cancelButton = this.rootElem.getElementsByClassName('btn-secondary')[0] as HTMLElement;
        cancelButton.classList.add('disabled');

        cancelButton.addEventListener('click', () => this.isCanceled = true);

        const calcButton = this.rootElem.getElementsByClassName('calc-weights')[0] as HTMLElement;
        calcButton.addEventListener('click', async () => {
            cancelButton.classList.remove('disabled');
            this.isCanceled = false;
            const previousContents = calcButton.innerHTML;
            calcButton.classList.add('disabled');
            calcButton.style.width = `${calcButton.getBoundingClientRect().width.toFixed(3)}px`;
            calcButton.innerHTML = `<i class="fa fa-spinner fa-spin"></i>&nbsp;Running`;
            this.container.scrollTo({top: 0});
            this.container.classList.add('pending');
            this.resultsViewer.setPending();
            const iterations = this.simUI.sim.getIterations();
            this.simUI.sim.setIterations(TypedEvent.nextEventID(), 3000);


            // 4273 = Ulduar, 3 = 10 Man, 5 = 25 Man
            let allItemsOfRaidForSlotChest = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotChest, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);
            let allItemsOfRaidForSlotHead = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotHead, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);
            let allItemsOfRaidForSlotNeck = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotNeck, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);
            let allItemsOfRaidForSlotShoulder = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotShoulder, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);
            let allItemsOfRaidForSlotBack = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotBack, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);
            let allItemsOfRaidForSlotWrist = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotWrist, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);
            let allItemsOfRaidForSlotHands = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotHands, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);
            let allItemsOfRaidForSlotWaist = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotWaist, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);
            let allItemsOfRaidForSlotLegs = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotLegs, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);
            let allItemsOfRaidForSlotFeet = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotFeet, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);
            let allItemsOfRaidForSlotTrinket = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotTrinket1, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);
            let allItemsOfRaidForSlotFinger = this.getAllItemsOfRaidForSlot(this.simUI, ItemSlot.ItemSlotFinger1, this.statsType.startsWith("U") ? 4273 : 0, this.statsType.endsWith("25") ? 5 : 3);

            this.totalSims = allItemsOfRaidForSlotBack.length + allItemsOfRaidForSlotChest.length + allItemsOfRaidForSlotHead.length +
                allItemsOfRaidForSlotNeck.length + allItemsOfRaidForSlotShoulder.length + allItemsOfRaidForSlotWrist.length +
                allItemsOfRaidForSlotHands.length + allItemsOfRaidForSlotWaist.length + allItemsOfRaidForSlotLegs.length +
                allItemsOfRaidForSlotFeet.length + allItemsOfRaidForSlotTrinket.length *2 + allItemsOfRaidForSlotFinger.length*2 + 1;

            let simResult = await this.simUI.sim.runRaidSim(TypedEvent.nextEventID(), (progress: ProgressMetrics) => {
                this.setSimProgress(progress);
            });
            let baseValue = simResult.raidMetrics.dps.avg;

            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotChest, allItemsOfRaidForSlotChest);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotHead, allItemsOfRaidForSlotHead);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotNeck, allItemsOfRaidForSlotNeck);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotShoulder, allItemsOfRaidForSlotShoulder);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotBack, allItemsOfRaidForSlotBack);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotWrist, allItemsOfRaidForSlotWrist);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotHands, allItemsOfRaidForSlotHands);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotWaist, allItemsOfRaidForSlotWaist);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotLegs,allItemsOfRaidForSlotLegs);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotFeet,allItemsOfRaidForSlotFeet);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotTrinket1,allItemsOfRaidForSlotTrinket);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotTrinket2,allItemsOfRaidForSlotTrinket);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotFinger1,allItemsOfRaidForSlotFinger);
            await this.getDpsDiffs(baseValue, resultList, ItemSlot.ItemSlotFinger2,allItemsOfRaidForSlotFinger);

            this.container.classList.remove('pending');
            this.resultsViewer.hideAll();
            calcButton.innerHTML = previousContents;
            calcButton.classList.remove('disabled');
            this.simUI.prevEpIterations = iterations;
            this.showResults(resultList);
        });

        // Default Gem Options
        const defaultGemDiv = document.createElement("div");
        if (this.autoGem) {
            defaultGemDiv.style.display = "flex";
        } else {
            defaultGemDiv.style.display = "none";
        }

        defaultGemDiv.classList.add("default-gem-container");
        const gemLabel = document.createElement("label");
        gemLabel.innerText = "Defaults for Auto Gem";
        defaultGemDiv.appendChild(gemLabel);

        const gemSocketsDiv = document.createElement("div");
        gemSocketsDiv.classList.add("sockets-container");

        Array<GemColor>(GemColor.GemColorRed, GemColor.GemColorYellow, GemColor.GemColorBlue, GemColor.GemColorMeta).forEach((socketColor, socketIndex) => {
            let gemFragment = document.createElement('fragment');
            gemFragment.innerHTML = `
          <div class="gem-socket-container">
            <img class="gem-icon" />
            <img class="socket-icon" />
          </div>
        `;

            const gemContainer = gemFragment.children[0] as HTMLElement;
            const gemIconElem = gemContainer.querySelector('.gem-icon') as HTMLImageElement;
            const socketIconElem = gemContainer.querySelector('.socket-icon') as HTMLImageElement;
            socketIconElem.src = getEmptyGemSocketIconUrl(socketColor);

            let selector: GemSelectorModal;

            let handleChoose = (itemData: ItemData<UIGem>) => {
                this.defaultGems[socketIndex] = itemData.item;
                ActionId.fromItemId(itemData.id).fill().then(filledId => {
                    gemIconElem.src = filledId.iconUrl;
                });
                selector.close();
            };

            let openGemSelector = () => {
                return () => {
                    if (selector == null) {
                        selector = new GemSelectorModal(this.simUI.rootElem, this.simUI, socketColor, handleChoose);
                    }
                    selector.show();
                }
            }

            gemIconElem.addEventListener("click", openGemSelector());
            gemContainer.addEventListener("click", openGemSelector());
            gemSocketsDiv.appendChild(gemContainer);
        });
        defaultGemDiv.appendChild(gemSocketsDiv);
        new BooleanPicker<DroptimizerMenu>(this.containerGemPicker, this, {
            label: "Auto Gem",
            labelTooltip: "When checked droptimizer will fill any un-filled gem sockets with default gems.",
            changedEvent: () => this.itemsChangedEmitter,
            getValue: () => this.autoGem,
            setValue: (id: EventID, obj: DroptimizerMenu, value: boolean) => {
                obj.autoGem = value
                if (value) {
                    defaultGemDiv.style.display = "flex";
                } else {
                    defaultGemDiv.style.display = "none";
                }
            }
        });

        this.containerGemPicker.appendChild(defaultGemDiv);
    }


    private async getDpsDiffs(baseValue: number, items: Result[], itemSlot: ItemSlot, allItemsOfRaidForSlot: UIItem[]) {
        if (this.isCanceled) {
            return;
        }

        let itemBefore = this.simUI.player.getEquippedItem(itemSlot);
        if (itemBefore == undefined)
            return;
        for (const itemToEquip of allItemsOfRaidForSlot) {
            let newGems = new Array<UIGem>();
            for (let i = 0; i < itemToEquip.gemSockets.length; i++){
                const newGem = itemToEquip.gemSockets[i];
                let gem = this.defaultGems[this.getSocketNumberToFind(newGem)];
                if (gem === undefined) {
                    continue;
                }
                newGems.push(gem)
            }
            let equip = new EquippedItem(itemToEquip, itemBefore?._enchant, newGems);
            this.simUI.player.equipItem(TypedEvent.nextEventID(), itemSlot, equip);

            this.currentSim++;
            let simResultNext = await this.simUI.sim.runRaidSim(TypedEvent.nextEventID(), (progress: ProgressMetrics) => {
                this.setSimProgress(progress);
            });

            let dpsResult = simResultNext.raidMetrics.dps.avg - baseValue;
            if (dpsResult > 0) {
                items.push(new Result(dpsResult, itemToEquip, itemBefore));
            }
        }
        this.simUI.player.equipItem(TypedEvent.nextEventID(), itemSlot, itemBefore);
    }


    private getSocketNumberToFind(color: GemColor) : number {
        switch (color){
            case GemColor.GemColorMeta:
                return 3;
            case GemColor.GemColorBlue:
                return 2;
            case GemColor.GemColorRed:
                return 0;
            case GemColor.GemColorYellow:
                return 1;
        }
        return 0;
    }
    private getAllItemsOfRaidForSlot(simUI: IndividualSimUI<any>, itemSlot: ItemSlot, raidId: number, difficulty: number) : UIItem[] {
        let items = simUI.player.getItems(itemSlot).filter(e => e.sources.some(itemSource => itemSource.source.oneofKind === "drop" && itemSource.source.drop.zoneId === raidId && itemSource.source.drop.difficulty === difficulty));

        items = items.filter(e => this.simUI.player.getEquippedItem(itemSlot)?.item.id !== e.id);
        let offSlot: ItemSlot | null = null;

        switch (itemSlot){
            case ItemSlot.ItemSlotFinger1:
                offSlot = ItemSlot.ItemSlotFinger2;
                break;
            case ItemSlot.ItemSlotFinger2:
                offSlot = ItemSlot.ItemSlotFinger1;
                break;
            case ItemSlot.ItemSlotTrinket1:
                offSlot = ItemSlot.ItemSlotTrinket2;
                break;
            case ItemSlot.ItemSlotTrinket2:
                offSlot = ItemSlot.ItemSlotTrinket1;
                break;
        }

        if (offSlot !== null) {
            let offSlotNonNull = (<ItemSlot>offSlot);
            items = items.filter(e => this.simUI.player.getEquippedItem(offSlotNonNull)?.item.id !== e.id);
        }
        return items;
    }

    private setSimProgress(progress: ProgressMetrics) {
        this.resultsViewer.setContent(`
			<div class="results-sim">
				<div class=""> ${this.currentSim} / ${this.totalSims}<br>simulations complete</div>
				<div class="">
					${progress.completedIterations} / ${progress.totalIterations}<br>iterations complete
				</div>
			</div>
		`);
    }


    showResults(resultList: Array<Result>) {
        resultList = resultList.sort((first, second) => second.dpsDiff - first.dpsDiff);
        this.tableBody.innerHTML = ``;
        for (const result of resultList){
            let row = this.makeTableRow(result);
            this.tableBody.appendChild(row);
        }
    }
    private makeTableRow(result: Result): HTMLElement {
        const row = document.createElement('tr');
        const makeWeightAndEpCellHtml = (result: Result, className: string): string => {


            let template = document.createElement('template');
            template.innerHTML = `
				<td class="stdev-cell ${className} type-weight">
					<span class="results-avg">${result.dpsDiff.toFixed(2)}</span>
				</td>
			`;

            return template.innerHTML;
        };

        row.innerHTML = `
			<td>${result.item.name}</td>
			${makeWeightAndEpCellHtml(result, 'damage-metrics')}
			<td class="current-ep"></td>
		`;

        return row;
    }

}
class Result {
    dpsDiff: number;
    item: UIItem;
    baseItem: EquippedItem;
    constructor(dpsDiff: number, item: UIItem, baseItem: EquippedItem) {
        this.dpsDiff = dpsDiff;
        this.item = item;
        this.baseItem = baseItem;
    }
}