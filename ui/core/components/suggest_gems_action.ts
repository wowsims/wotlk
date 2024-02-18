import { IndividualSimUI } from '../individual_sim_ui.js';
import { Player } from '../player.js';
import { Sim } from '../sim.js';
import { Gear } from '../proto_utils/gear.js';
import { EquippedItem } from '../proto_utils/equipped_item.js';
import { TypedEvent } from '../typed_event.js';
import { Stats } from '../proto_utils/stats.js';
import { GemColor, Stat, Profession, ItemSlot, Spec } from '../proto/common.js';

interface GemCapsData {
	gemId: number
	statCaps: Stats
}

interface SocketData {
	itemSlot: ItemSlot
	socketIdx: number
}

interface SocketBonusData {
	itemSlot: ItemSlot | null
	socketBonus: number
}

abstract class GemOptimizer {
	protected readonly player: Player<Spec>;
	protected readonly sim: Sim;
	protected readonly gemPriorityByColor: Record<GemColor, Array<GemCapsData>>; 
	abstract metaGemID: number;	
	static allGemColors: Array<GemColor> = [GemColor.GemColorRed, GemColor.GemColorYellow, GemColor.GemColorBlue];
	epWeights!: Stats;
	useJcGems!: boolean;
	isBlacksmithing!: boolean;
	numSocketedJcGems!: number;
	jcUpgradePriority: Array<GemCapsData>;

	static jcUpgradesById: Record<number, number> = {
		40118: 42154,
		40125: 42156,
		40112: 42143,
		40111: 42142,
		40119: 36767,
	};

	constructor(simUI: IndividualSimUI<any>) {
		this.player = simUI.player;
		this.sim = simUI.sim;

		// Initialize empty arrays of gem priorities for each socket color
		this.gemPriorityByColor = {} as Record<GemColor, Array<GemCapsData>>;
		
		for (var gemColor of GemOptimizer.allGemColors) {
			this.gemPriorityByColor[gemColor] = new Array<GemCapsData>();
		}

		this.jcUpgradePriority = new Array<GemCapsData>();	

		simUI.addAction('Suggest Gems', 'suggest-gems-action', async () => {
			this.optimizeGems();
		});
	}

	async optimizeGems() {
		// First, clear all existing gems
		let optimizedGear = this.player.getGear().withoutGems();
		this.numSocketedJcGems = 0;

		// Store relevant player attributes for use in optimizations
		this.epWeights = this.player.getEpWeights();
		this.useJcGems = this.player.hasProfession(Profession.Jewelcrafting);
		this.isBlacksmithing = this.player.isBlacksmithing();

		/*
		 * Use subclass-specific logic to rank order gems of each color by value
		 * and calculate the associated stat caps for each gem (when applicable).
		 */
		const ungemmedStats = await this.updateGear(optimizedGear);
		this.updateGemPriority(optimizedGear, ungemmedStats);

		// Next, socket and activate the meta
		optimizedGear = optimizedGear.withMetaGem(this.sim.db.lookupGem(this.metaGemID));
		optimizedGear = this.activateMetaGem(optimizedGear);
		await this.updateGear(optimizedGear);

		// Now loop through all gem colors where a priority list has been defined
		for (var gemColor of GemOptimizer.allGemColors) {
			if (this.gemPriorityByColor[gemColor].length > 0) {
				optimizedGear = await this.fillGemsByColor(optimizedGear, gemColor);

				// Also substitute JC gems by priority while respecting stat caps
				if (this.useJcGems) {
					optimizedGear = await this.substituteJcGems(optimizedGear);
				}
			}
		}
	}
	
	async updateGear(gear: Gear): Promise<Stats> {
		this.player.setGear(TypedEvent.nextEventID(), gear);
		await this.sim.updateCharacterStats(TypedEvent.nextEventID());
		return Stats.fromProto(this.player.getCurrentStats().finalStats);
	}

	/**
	 * Helper method for meta gem activation.
	 *
	 * @remarks
	 * Based on the ansatz that most specs are forced to use a suboptimal gem color in
	 * order to statisfy their meta requirements. As a result, it is helpful to
	 * compute the item slot in a gear set that provides the strongest socket bonus 
	 * for that color, since this should minimize the "cost" of activation.
	 *
	 * @param gear - Ungemmed gear set
	 * @param color - Socket color used for meta gem activation
	 * @param singleOnly - If true, exclude items containing more than one socket of the specified color. If false, instead normalize the socket bonus by the number of such sockets.
	 * @param blacklistedColor - If non-null, exclude items containing any sockets of this color (assumed to be different from the color used for activation).
	 * @returns Optimal item slot for activation under the specified constraints, or null if not found.
	 */	
	findStrongestSocketBonus(gear: Gear, color: GemColor, singleOnly: boolean, blacklistedColor: GemColor | null): SocketBonusData {
		let optimalSlot: ItemSlot | null = null;
		let maxSocketBonusEP: number = 1e-8;

		for (var slot of gear.getItemSlots()) {
			const item = gear.getEquippedItem(slot);

			if (!item) {
				continue;
			}

			if (item.numSocketsOfColor(blacklistedColor) != 0) {
				continue;
			}

			const numSockets = item.numSocketsOfColor(color);

			if ((numSockets == 0) || (singleOnly && (numSockets != 1))) {
				continue;
			}

			const socketBonusEP = new Stats(item.item.socketBonus).computeEP(this.epWeights);
			const normalizedEP = socketBonusEP / numSockets;

			if (normalizedEP > maxSocketBonusEP) {
				optimalSlot = slot;
				maxSocketBonusEP = normalizedEP;
			}
		}

		return { itemSlot: optimalSlot, socketBonus: maxSocketBonusEP };
	}
	
	socketGemInFirstMatchingSocket(gear: Gear, itemSlot: ItemSlot | null, colorToMatch: GemColor, gemId: number): Gear {
		if (itemSlot != null) {
			const item = gear.getEquippedItem(itemSlot);

			if (!item) {
				return gear;
			}

			for (const [socketIdx, socketColor] of item!.allSocketColors().entries()) {
				if (socketColor == colorToMatch) {
					return gear.withEquippedItem(itemSlot, item!.withGem(this.sim.db.lookupGem(gemId), socketIdx), true);
				}
			}
		}

		return gear;
	}

	async fillGemsByColor(gear: Gear, color: GemColor): Promise<Gear> {
		const socketList = this.findSocketsByColor(gear, color);
		return await this.fillGemsToCaps(gear, socketList, this.gemPriorityByColor[color], 0, 0);
	}
	
	/**
	 * Shared wrapper for compiling eligible sockets for each gem priority list.
	 *
	 * @remarks
	 * Subclasses are required to implement the allowGemInSocket method, which
	 * contains the (spec-specific) logic on when to match socket bonuses etc.
	 *
	 * @param gear - Partially gemmed gear set
	 * @param color - Color associated with a single gem priority list
	 * @returns Array of sockets that will be filled using the priority list associated with the specified color.
	 */	
	findSocketsByColor(gear: Gear, color: GemColor): Array<SocketData> {
		const socketList = new Array<SocketData>();

		for (var slot of gear.getItemSlots()) {
			const item = gear.getEquippedItem(slot);

			if (!item) {
				continue;
			}

			for (const [socketIdx, socketColor] of item.curSocketColors(this.isBlacksmithing).entries()) {
				if (item!.hasSocketedGem(socketIdx)) {
					continue;
				}

				if (this.allowGemInSocket(color, socketColor, slot, item)) {
					socketList.push({ itemSlot: slot, socketIdx: socketIdx });
				}
			}
		}

		return socketList;
	}

	async substituteJcGems(gear: Gear): Promise<Gear> {
		let updatedGear: Gear = gear;
		let gemIdx = 0;

		while ((this.numSocketedJcGems < 3) && (gemIdx < this.jcUpgradePriority.length)) {
			const gemData = this.jcUpgradePriority[gemIdx];
			const baseGem = this.sim.db.lookupGem(gemData.gemId);

			if (!updatedGear.getAllGems(this.isBlacksmithing).includes(baseGem!)) {
				gemIdx += 1;
				continue;
			}

			const upgradedGem = this.sim.db.lookupGem(GemOptimizer.jcUpgradesById[gemData.gemId]);
			const testGear = updatedGear.withSingleGemSubstitution(baseGem, upgradedGem, this.isBlacksmithing);
			const newStats = await this.updateGear(testGear);

			if (newStats.belowCaps(gemData.statCaps)) {
				updatedGear = testGear;
				this.numSocketedJcGems += 1;
			} else {
				await this.updateGear(updatedGear);
				gemIdx += 1;
			}	
		}

		return updatedGear;
	}
	
	async fillGemsToCaps(gear: Gear, socketList: Array<SocketData>, gemCaps: Array<GemCapsData>, numPasses: number, firstIdx: number): Promise<Gear> {
		let updatedGear: Gear = gear;
		const currentGem = this.sim.db.lookupGem(gemCaps[numPasses].gemId);

		// On the first pass, we simply fill all sockets with the highest priority gem
		if (numPasses == 0) {
			for (var socketData of socketList.slice(firstIdx)) {
				updatedGear = updatedGear.withGem(socketData.itemSlot, socketData.socketIdx, currentGem);
			}
		}

		// If we are below the relevant stat cap for the gem we just filled on the last pass, then we are finished.
		let newStats = await this.updateGear(updatedGear);
		const currentCap = gemCaps[numPasses].statCaps;

		if (newStats.belowCaps(currentCap) || (numPasses == gemCaps.length - 1)) {
			return updatedGear;
		}

		// If we exceeded the stat cap, then work backwards through the socket list and replace each gem with the next highest priority option until we are below the cap
		const nextGem = this.sim.db.lookupGem(gemCaps[numPasses + 1].gemId);
		const nextCap = gemCaps[numPasses + 1].statCaps;
		let capForReplacement = currentCap.subtract(nextCap);

		if (currentCap.computeEP(capForReplacement) <= 0) {
			capForReplacement = currentCap;
		}

		for (var idx = socketList.length - 1; idx >= firstIdx; idx--) {
			if (newStats.belowCaps(capForReplacement)) {
				break;
			}

			updatedGear = updatedGear.withGem(socketList[idx].itemSlot, socketList[idx].socketIdx, nextGem);
			newStats = await this.updateGear(updatedGear);
		}

		// Now run a new pass to check whether we've exceeded the next stat cap
		let nextIdx = idx + 1;

		if (!newStats.belowCaps(currentCap)) {
			nextIdx = firstIdx;
		}

		return await this.fillGemsToCaps(updatedGear, socketList, gemCaps, numPasses + 1, nextIdx);
	}

	abstract activateMetaGem(gear: Gear): Gear;

	abstract updateGemPriority(ungemmedGear: Gear, passiveStats: Stats): void;

	abstract allowGemInSocket(gemColor: GemColor, socketColor: GemColor, itemSlot: ItemSlot, item: EquippedItem): boolean;
}

export class PhysicalDPSGemOptimizer extends GemOptimizer {
	metaGemID: number = 41398; // Relentless Earthsiege Diamond
	arpSlop: number = 11;
	expSlop: number = 4;
	hitTarget: number = 8. * 32.79;
	hitSlop: number = 4;
	useArpGems: boolean;
	useExpGems: boolean;
	useAgiGems: boolean;
	useStrGems: boolean;
	arpTarget!: number;
	passiveArp!: number;
	arpStackDetected!: boolean;
	passiveHit!: number;
	tearSlot!: ItemSlot | null;

	constructor(simUI: IndividualSimUI<any>, useArpGems: boolean, useExpGems: boolean, useAgiGems: boolean, useStrGems: boolean) {
		super(simUI);
		this.useArpGems = useArpGems;
		this.useExpGems = useExpGems;
		this.useAgiGems = useAgiGems;
		this.useStrGems = useStrGems;
	}

	updateGemPriority(ungemmedGear: Gear, passiveStats: Stats) {
		// First calculate any gear-dependent stat caps.
		this.arpTarget = this.calcArpTarget(ungemmedGear);
		const critCap = this.calcCritCap(ungemmedGear);
		const expCap = new Stats().withStat(Stat.StatExpertise, this.calcExpTarget() + this.expSlop);
		this.passiveHit = passiveStats.getStat(Stat.StatMeleeHit);
		const hitCap = new Stats().withStat(Stat.StatMeleeHit, this.hitTarget + this.hitSlop);

		// Reset optimal Tear slot from prior calculations
		this.tearSlot = null;
		
		/*
		 * For specs that gem ArP, determine whether the current gear
		 * configuration will optimally hard stack Fractured gems or not.
		 */
		this.passiveArp = passiveStats.getStat(Stat.StatArmorPenetration);
		this.arpStackDetected = this.detectArpStackConfiguration(ungemmedGear);

		/*
		 * Use tighter constraint on overcapping ArP for hard stack setups, so as
		 * to reduce the number of missed yellow socket bonuses.
		 */
		const arpSlop = this.arpStackDetected ? 4 : this.arpSlop;
		const arpCap = new Stats().withStat(Stat.StatArmorPenetration, this.arpTarget + arpSlop);

		// Update red gem priority
		const redGemCaps = new Array<GemCapsData>();

		// Fractured Cardinal Ruby
		if (this.useArpGems) {
			redGemCaps.push({ gemId: 40117, statCaps: arpCap });
		}

		// Precise Cardinal Ruby
		if (this.useExpGems) {
			redGemCaps.push({ gemId: 40118, statCaps: expCap });
		}

		// Delicate Cardinal Ruby
		if (this.useAgiGems) {
			redGemCaps.push({ gemId: 40112, statCaps: critCap });
		}

		// Bold Cardinal Ruby
		if (this.useStrGems) {
			redGemCaps.push({ gemId: 40111, statCaps: new Stats() });
		}

		this.gemPriorityByColor[GemColor.GemColorRed] = redGemCaps;

		// Update yellow gem priority
		const yellowGemCaps = new Array<GemCapsData>();

		// Accurate Ametrine
		if (this.useExpGems) {
			yellowGemCaps.push({ gemId: 40162, statCaps: hitCap.add(expCap) });
		}

		// Rigid Ametrine
		yellowGemCaps.push({ gemId: 40125, statCaps: hitCap });

		// Fractured Cardinal Ruby
		if (this.arpStackDetected) {
			yellowGemCaps.push({ gemId: 40117, statCaps: arpCap });
		}
		
		// Accurate Ametrine (needed to add twice to catch some edge cases)
		if (this.useExpGems) {
			yellowGemCaps.push({ gemId: 40162, statCaps: hitCap.add(expCap) });
		}

		// Glinting Ametrine
		if (this.useAgiGems) {
			yellowGemCaps.push({ gemId: 40148, statCaps: hitCap.add(critCap) });
		}

		// Etched Ametrine
		if (this.useStrGems) {
			yellowGemCaps.push({ gemId: 40143, statCaps: hitCap });
		}

		// Deadly Ametrine
		if (this.useAgiGems) {
			yellowGemCaps.push({ gemId: 40147, statCaps: critCap });
		}

		// Inscribed Ametrine
		if (this.useStrGems) {
			yellowGemCaps.push({ gemId: 40142, statCaps: critCap });
		}

		// Fierce Ametrine
		if (this.useStrGems) {
			yellowGemCaps.push({ gemId: 40146, statCaps: new Stats() });
		}
		
		this.gemPriorityByColor[GemColor.GemColorYellow] = yellowGemCaps;

		// Update JC upgrade priority
		this.jcUpgradePriority = new Array<GemCapsData>();
		
		if (this.useExpGems) {
			this.jcUpgradePriority.push({ gemId: 40118, statCaps: expCap });
		}
		
		if (this.useAgiGems) {
			this.jcUpgradePriority.push({ gemId: 40112, statCaps: critCap });
		}

		if (this.useStrGems) {
			this.jcUpgradePriority.push({ gemId: 40111, statCaps: new Stats() });
		}
	}

	detectArpStackConfiguration(ungemmedGear: Gear): boolean {
		if (!this.useArpGems) {
			return false;
		}

		/*
		 * Generate a "dummy" list of red sockets in order to determine whether
		 * ignoring yellow socket bonuses to stack more ArP gems will be correct.
		 * Subtract 2 from the length of this list to account for meta gem +
		 * Nightmare Tear.
		 */
		const dummyRedSocketList = this.findSocketsByColor(ungemmedGear, GemColor.GemColorRed);
		const numRedSockets = dummyRedSocketList.length - 2;
		let projectedArp = this.passiveArp + 20 * numRedSockets;

		if (this.useJcGems) {
			projectedArp += 42;
		}

		return (this.arpTarget > 1000) && (projectedArp > 648) && (projectedArp + 20 < this.arpTarget + 4);
	}

	activateMetaGem(gear: Gear): Gear {
		/*
		 * Use a single Nightmare Tear for meta activation. Prioritize blue
		 * sockets for it if possible, and fall back to yellow sockets if not.
		 */
		const blueSlotCandidate = this.findBlueTearSlot(gear);
		const yellowSlotCandidate = this.findYellowTearSlot(gear);

		let tearColor = GemColor.GemColorBlue;
		this.tearSlot = blueSlotCandidate.itemSlot;

		if ((this.tearSlot == null) || (this.arpStackDetected && (yellowSlotCandidate.socketBonus > blueSlotCandidate.socketBonus))) {
			tearColor = GemColor.GemColorYellow;
			this.tearSlot = yellowSlotCandidate.itemSlot;
		}

		return this.socketTear(gear, tearColor);
	}
	
	socketTear(gear: Gear, tearColor: GemColor): Gear {
		return this.socketGemInFirstMatchingSocket(gear, this.tearSlot, tearColor, 49110);
	}
	
	findBlueTearSlot(gear: Gear): SocketBonusData {
		// Eligible Tear slots have only one blue socket max.
		const singleOnly = true;

		/*
		 * Additionally, for hard ArP stack configurations, only use blue sockets
		 * for Tear if there are no yellow sockets in that item slot, since hard
		 * ArP stacks ignore yellow socket bonuses in favor of stacking more
		 * Fractured gems.
		 */
		const blacklistedColor = this.arpStackDetected ? GemColor.GemColorYellow : null;

		return this.findStrongestSocketBonus(gear, GemColor.GemColorBlue, singleOnly, blacklistedColor);
	}

	findYellowTearSlot(gear: Gear): SocketBonusData {
		return this.findStrongestSocketBonus(gear, GemColor.GemColorYellow, false, GemColor.GemColorBlue);
	}
	
	allowGemInSocket(gemColor: GemColor, socketColor: GemColor, itemSlot: ItemSlot, item: EquippedItem): boolean {
		const ignoreYellowSockets = ((item!.numSocketsOfColor(GemColor.GemColorBlue) > 0) && (itemSlot != this.tearSlot));
		let matchYellowSocket = false;
		
		if ((socketColor == GemColor.GemColorYellow) && !ignoreYellowSockets) {
			matchYellowSocket = new Stats(item.item.socketBonus).computeEP(this.epWeights) > 1e-8;
		}

		return ((gemColor == GemColor.GemColorYellow) && matchYellowSocket) || ((gemColor == GemColor.GemColorRed) && !matchYellowSocket);
	}
	
	findSocketsByColor(gear: Gear, color: GemColor): Array<SocketData> {
		const socketList = super.findSocketsByColor(gear, color);

		if (this.arpStackDetected && (color == GemColor.GemColorYellow)) {
			this.sortYellowSockets(gear, socketList);
		}

		return socketList;
	}
	
	sortYellowSockets(gear: Gear, yellowSocketList: Array<SocketData>) {
		yellowSocketList.sort((a,b) => {
			// If both yellow sockets belong to the same item, then treat them equally.
			const slot1 = a.itemSlot;
			const slot2 = b.itemSlot;

			if (slot1 == slot2) {
				return 0;
			}

			// If an item already has a Nightmare Tear socketed, then bump up any yellow sockets in it to highest priority.
			if (slot1 == this.tearSlot) {
				return -1;
			}

			if (slot2 == this.tearSlot) {
				return 1;
			}

			// For all other cases, sort by the ratio of the socket bonus value divided by the number of yellow sockets required to activate it.
			const item1 = gear.getEquippedItem(slot1);
			const bonus1 = new Stats(item1!.item.socketBonus).computeEP(this.epWeights);
			const item2 = gear.getEquippedItem(slot2);
			const bonus2 = new Stats(item2!.item.socketBonus).computeEP(this.epWeights);
			return bonus2 / item2!.numSocketsOfColor(GemColor.GemColorYellow) - bonus1 / item1!.numSocketsOfColor(GemColor.GemColorYellow);
		});
	}
	
	calcArpTarget(gear: Gear): number {
		let arpTarget = 1399;

		/*
		 * First handle ArP proc trinkets. If more than one of these are equipped
		 * simultaneously, it is assumed that the user is desyncing them via ICD
		 * resets, such that the soft cap is set by the strongest proc.
		 */
		if (gear.hasTrinket(45931)) {
			arpTarget -= 751; // Mjolnir Runestone
		} else if (gear.hasTrinket(50198)) {
			arpTarget -= 678; // Needle-Encrusted Scorpion
		} else if (gear.hasTrinket(40256)) {
			arpTarget -= 612; // Grim Toll
		}

		// Then check for Executioner enchant
		const weapon = gear.getEquippedItem(ItemSlot.ItemSlotMainHand);

		if (weapon?.enchant?.effectId == 3225) {
			arpTarget -= 120;
		}

		return arpTarget;
	}

	calcExpTarget(): number {
		return 6.5 * 32.79;
	}
	
	calcCritCap(gear: Gear): Stats {
		/*
		 * Only some specs incorporate Crit soft caps into their gemming logic, so
		 * the parent method here simply returns an empty Stats object (meaning
		 * that Crit cap will just be ignored elsewhere in the code). Custom
		 * spec-specific subclasses can override this as desired.
		 */
		return new Stats();
	}
	
	async fillGemsByColor(gear: Gear, color: GemColor): Promise<Gear> {
		/*
		 * Parent logic substitutes JC gems after filling normal gems first, but
		 * for specs that gem ArP, it is more optimal to pre-fill some Fractured
		 * Dragon's Eyes if doing so gets us closer to the target.
		 */
		let updatedGear: Gear = gear;

		if ((color == GemColor.GemColorRed) && this.useArpGems && this.useJcGems) {
			updatedGear = this.optimizeJcArpGems(updatedGear);
		}

		// Likewise, if we still have JC gems available after finishing the red gems, then force utilization of JC Hit gems if possible.
		if ((color == GemColor.GemColorYellow) && this.useJcGems && (this.numSocketedJcGems < 3)) {
			updatedGear = this.fillJcHitGems(updatedGear);
		}

		return await super.fillGemsByColor(updatedGear, color);
	}

	calcDistanceToArpTarget(numJcArpGems: number, numRedSockets: number): number {
		const numNormalArpGems = Math.max(0, Math.min(numRedSockets - 3, Math.floor((this.arpTarget + this.arpSlop - this.passiveArp - 34 * numJcArpGems) / 20)));
		const projectedArp = this.passiveArp + 34 * numJcArpGems + 20 * numNormalArpGems;
		return Math.abs(projectedArp - this.arpTarget);
	}

	optimizeJcArpGems(gear: Gear): Gear {
		// First determine how many of the JC gems should be 34 ArP gems
		const redSocketList = this.findSocketsByColor(gear, GemColor.GemColorRed);
		const numRedSockets = redSocketList.length;
		let optimalJcArpGems = [0,1,2,3].reduce((m,x)=> this.calcDistanceToArpTarget(m, numRedSockets)<this.calcDistanceToArpTarget(x, numRedSockets) ? m:x);
		optimalJcArpGems = Math.min(optimalJcArpGems, numRedSockets);

		// Now socket just those gems, saving other JC substitutions for later
		let updatedGear: Gear = gear;

		for (let i = 0; i < optimalJcArpGems; i++) {
			updatedGear = updatedGear.withGem(redSocketList[i].itemSlot, redSocketList[i].socketIdx, this.sim.db.lookupGem(42153));
		}

		this.numSocketedJcGems = optimalJcArpGems;
		return updatedGear;
	}

	fillJcHitGems(gear: Gear): Gear {
		const yellowSocketList = this.findSocketsByColor(gear, GemColor.GemColorYellow);
		const maxJcHitGems = Math.min(3 - this.numSocketedJcGems, yellowSocketList.length);
		const desiredJcHitGems = Math.max(0, Math.floor((this.hitTarget + this.hitSlop - this.passiveHit) / 34));
		const numJcHitGems = Math.min(desiredJcHitGems, maxJcHitGems);

		let updatedGear: Gear = gear;

		for (let i = 0; i < numJcHitGems; i++) {
			updatedGear = updatedGear.withGem(yellowSocketList[i].itemSlot, yellowSocketList[i].socketIdx, this.sim.db.lookupGem(42156));
		}

		this.numSocketedJcGems += numJcHitGems;
		return updatedGear;
	}
}

export class TankGemOptimizer extends GemOptimizer {
	metaGemID: number = 41380; // Austere Earthsiege Diamond
	
	updateGemPriority(ungemmedGear: Gear, passiveStats: Stats) {
		// Base class just stuffs pure Stamina gems everywhere
		const blueGemCaps = new Array<GemCapsData>();
		blueGemCaps.push({ gemId: 40119, statCaps: new Stats() });
		this.gemPriorityByColor[GemColor.GemColorBlue] = blueGemCaps;
		this.jcUpgradePriority = blueGemCaps;
	}
	
	activateMetaGem(gear: Gear): Gear {
		/*
		 * Use a single Shifting Dreadstone gem for meta activation, in the slot
		 * with the strongest bonus for a single red socket.
		 */
		return this.socketGemInFirstMatchingSocket(gear, this.findStrongestSocketBonus(gear, GemColor.GemColorRed, true, GemColor.GemColorYellow).itemSlot, GemColor.GemColorRed, 40130);
	}
	
	allowGemInSocket(gemColor: GemColor, socketColor: GemColor, itemSlot: ItemSlot, item: EquippedItem): boolean {
		return gemColor == GemColor.GemColorBlue;
	}
}
