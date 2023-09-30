import { Player } from "../core/player";
import { GemColor, ItemSlot, Profession, Spec, Stat } from "../core/proto/common";
import { Gear } from "../core/proto_utils/gear";
import { Stats } from "../core/proto_utils/stats";
import { Sim } from "../core/sim";
import { TypedEvent } from "../core/typed_event";
import * as Mechanics from '../core/constants/mechanics.js';

/***
 * WARNING: Currently only optimised for Arp/Exp/Hit gemming the following specs;
 * - Feral
 * - Warrior
 * - Hunter
 */
type AutoGemSpec = Spec.SpecWarrior | Spec.SpecFeralDruid | Spec.SpecHunter

enum GemsByStats {
  Str = 40111,
  Agi = 40112,
  Arp = 40117,
  Exp = 40118,
  Hit = 40125,
  Str_Crit = 40142,
  Str_Hit = 40143,
  Str_Haste = 40146,
  Agi_Crit = 40147,
  Agi_Hit = 40148,
  Exp_Hit = 40162,
}

/**
 * Add spec specific slop to the real ArP cap
 */
const calcArpCap = (arpTarget: number, player: Player<AutoGemSpec>): Stats => {
  let arpCap = arpTarget
  // Sets additional "slop" to allow for minor overcapping
  switch (player.spec) {
    case Spec.SpecFeralDruid:
      arpCap += 11
      break;
    default:
      arpCap += 4
      break;
  }
  return new Stats().withStat(Stat.StatArmorPenetration, arpCap);
}

/**
 * Calculate the real ArP cap value the player needs
 */
const calcArpTarget = (gear: Gear): Stats => {
  let arpCap = 1399;

  // Mjolnir Runestone
  if (gear.hasTrinket(45931)) {
    arpCap -= 751;
  }

  // Grim Toll
  if (gear.hasTrinket(40256)) {
    arpCap -= 612;
  }

  // Executioner enchant
  const weapon = gear.getEquippedItem(ItemSlot.ItemSlotMainHand);
  if (weapon?.enchant?.effectId === 3225) {
    arpCap -= 120;
  }

  return new Stats().withStat(Stat.StatArmorPenetration, arpCap);
}

/**
 * Calculate the value still needed to hit ArP cap
 */
const calcDistanceToArpTarget = (numJcArpGems: number, passiveArp: number, numRedSockets: number, arpCap: number, arpTarget: number): number => {
  const numNormalArpGems = Math.max(0, Math.min(numRedSockets - 3, Math.floor((arpCap - passiveArp - 34 * numJcArpGems) / 20)));
  const projectedArp = passiveArp + 34 * numJcArpGems + 20 * numNormalArpGems;
  return Math.abs(projectedArp - arpTarget);
}

/**
 * Calculate the Expertise cap value and add spec specific slop
 */
const calcExpCap = (player: Player<AutoGemSpec>): Stats => {
  let expCap = 6.5 * 32.79;
  const talents = player.getTalents()

  switch (player.spec) {
    case Spec.SpecWarrior:
      if ('weaponMastery' in talents) {
        const weaponMastery = talents.weaponMastery;
        const hasWeaponMasteryTalent = !!weaponMastery;

        if (hasWeaponMasteryTalent) {
          expCap -=
            weaponMastery * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION;
        }
      }
      break;
  }

  expCap += 4 // Add 4 as default slop

  return new Stats().withStat(Stat.StatExpertise, expCap);
}

/**
 * Calculate the Crit cap value
 */
const calcCritCap = (gear: Gear): Stats => {
  const baseCritCapPercentage = 75.8 + 3; // includes 3% Crit debuff
  let agiProcs = 0;

  if (gear.hasRelic(47668)) {
    agiProcs += 200;
  }

  if (gear.hasRelic(50456)) {
    agiProcs += 44 * 5;
  }

  if (gear.hasTrinket(47131) || gear.hasTrinket(47464)) {
    agiProcs += 510;
  }

  if (gear.hasTrinket(47115) || gear.hasTrinket(47303)) {
    agiProcs += 450;
  }

  if (gear.hasTrinket(44253) || gear.hasTrinket(42987)) {
    agiProcs += 300;
  }

  return new Stats().withStat(Stat.StatMeleeCrit, (baseCritCapPercentage - agiProcs * 1.1 * 1.06 * 1.02 / 83.33) * 45.91);
}

/**
 * Calculate the Crit value
 * and add spec specific slop
 */
const calcHitCap = (player: Player<AutoGemSpec>): Stats => {
  let hitCap = 8. * 32.79

  // Sets additional "slop" to allow for minor overcapping
  switch (player.spec) {
    default:
      hitCap += 11
      break;
  }

  return new Stats().withStat(Stat.StatMeleeHit, hitCap)
}

const optimizeJewelCraftingGems = (sim: Sim, player: Player<AutoGemSpec>, gear: Gear, redSocketList: [ItemSlot, number][], arpCap: Stats, arpTarget: number): Gear => {
  const passiveStats = Stats.fromProto(player.getCurrentStats().finalStats);
  const passiveArp = passiveStats.getStat(Stat.StatArmorPenetration);
  const arpCapValue = arpCap.getStat(Stat.StatArmorPenetration);
  const numRedSockets = redSocketList.length;
  const isBelowCritCap = passiveStats.belowCaps(calcCritCap(gear));

  // First determine how many of the JC gems should be 34 ArP gems
  let optimalJcArpGems = 0;
  let minDistanceToArpTarget = calcDistanceToArpTarget(0, passiveArp, numRedSockets, arpCapValue, arpTarget);

  for (let i = 1; i <= 3; i++) {
    const distanceToArpTarget = calcDistanceToArpTarget(i, passiveArp, numRedSockets, arpCapValue, arpTarget);

    if (distanceToArpTarget < minDistanceToArpTarget) {
      optimalJcArpGems = i;
      minDistanceToArpTarget = distanceToArpTarget;
    }
  }

  // Now actually socket the gems
  let updatedGear: Gear = gear;
  for (let i = 0; i < 3; i++) {
    let gemId: number | null = null;

    switch (player.spec) {
      case Spec.SpecHunter:
        gemId = 42143 // Agi by default
      case Spec.SpecWarrior:
        gemId = 42142 // Str by default
      case Spec.SpecFeralDruid:
        gemId = 42142 // Str by default
    }


    if (i < optimalJcArpGems) {
      gemId = 42153; // ArP
    } else if ((player.spec === Spec.SpecFeralDruid) && isBelowCritCap) {
      gemId = 42143; // Below crit swap to Agi
    }

    if (gemId) updatedGear = updatedGear.withGem(redSocketList[i][0], redSocketList[i][1], sim.db.lookupGem(gemId));
  }

  return updatedGear;
}

const fillGemsToCaps = async (sim: Sim, player: Player<AutoGemSpec>, gear: Gear, socketList: Array<[ItemSlot, number]>, gemCaps: Array<[number, Stats]>, numPasses: number, firstIdx: number): Promise<Gear> => {
  let updatedGear: Gear = gear;
  const currentGem = sim.db.lookupGem(gemCaps[numPasses][0]);

  // On the first pass, we simply fill all sockets with the highest priority gem
  if (numPasses === 0) {
    for (const [itemSlot, socketIdx] of socketList.slice(firstIdx)) {
      updatedGear = updatedGear.withGem(itemSlot, socketIdx, currentGem);
    }
  }

  // If we are below the relevant stat cap for the gem we just filled on the last pass, then we are finished.
  let newStats = await updateGear(sim, player, updatedGear);
  const currentCap = gemCaps[numPasses][1];

  if (newStats.belowCaps(currentCap) || (numPasses === gemCaps.length - 1)) {
    return updatedGear;
  }

  // If we exceeded the stat cap, then work backwards through the socket list and replace each gem with the next highest priority option until we are below the cap
  const nextGem = sim.db.lookupGem(gemCaps[numPasses + 1][0]);
  const nextCap = gemCaps[numPasses + 1][1];
  let capForReplacement = currentCap;

  switch (player.spec) {
    case Spec.SpecFeralDruid:
      capForReplacement = currentCap.subtract(nextCap);
      if (currentCap.computeEP(capForReplacement) <= 0) {
        capForReplacement = currentCap;
      }
      break;
    case Spec.SpecWarrior:
    case Spec.SpecHunter:
      if ((numPasses > 0) && !currentCap.equals(nextCap)) {
        capForReplacement = currentCap.subtract(nextCap);
      }
      break
  }

  for (var idx = socketList.length - 1; idx >= firstIdx; idx--) {
    if (newStats.belowCaps(capForReplacement)) {
      break;
    }

    const [itemSlot, socketIdx] = socketList[idx];
    updatedGear = updatedGear.withGem(itemSlot, socketIdx, nextGem);
    newStats = await updateGear(sim, player, updatedGear);
  }

  // Now run a new pass to check whether we've exceeded the next stat cap
  let nextIdx = idx + 1;

  if (!newStats.belowCaps(currentCap)) {
    nextIdx = firstIdx;
  }

  return await fillGemsToCaps(sim, player, updatedGear, socketList, gemCaps, numPasses + 1, nextIdx);
}

const updateGear = async (sim: Sim, player: Player<AutoGemSpec>, gear: Gear): Promise<Stats> => {
  player.setGear(TypedEvent.nextEventID(), gear);
  await sim.updateCharacterStats(TypedEvent.nextEventID());
  return Stats.fromProto(player.getCurrentStats().finalStats);
}
const findBlueTearSlot = (gear: Gear, epWeights: Stats): ItemSlot | null => {
  let tearSlot: ItemSlot | null = null;
  let maxBlueSocketBonusEP: number = 1e-8;

  for (var slot of gear.getItemSlots()) {
    const item = gear.getEquippedItem(slot);

    if (!item) {
      continue;
    }

    if (item.numSocketsOfColor(GemColor.GemColorBlue) !== 1) {
      continue;
    }

    const socketBonusEP = new Stats(item.item.socketBonus).computeEP(epWeights);

    if (socketBonusEP > maxBlueSocketBonusEP) {
      tearSlot = slot;
      maxBlueSocketBonusEP = socketBonusEP;
    }
  }

  return tearSlot;
}

const findYellowTearSlot = (gear: Gear, epWeights: Stats): ItemSlot | null => {
  let tearSlot: ItemSlot | null = null;
  let maxYellowSocketBonusEP: number = 1e-8;

  for (var slot of gear.getItemSlots()) {
    const item = gear.getEquippedItem(slot);

    if (!item) {
      continue;
    }

    if (item.numSocketsOfColor(GemColor.GemColorBlue) !== 0) {
      continue;
    }

    const numYellowSockets = item!.numSocketsOfColor(GemColor.GemColorYellow);

    if (numYellowSockets === 0) {
      continue;
    }

    const socketBonusEP = new Stats(item.item.socketBonus).computeEP(epWeights);
    const normalizedEP = socketBonusEP / numYellowSockets;

    if (normalizedEP > maxYellowSocketBonusEP) {
      tearSlot = slot;
      maxYellowSocketBonusEP = normalizedEP;
    }
  }

  return tearSlot;
}

const socketTear = (sim: Sim, gear: Gear, tearSlot: ItemSlot | null, tearColor: GemColor): Gear => {
  if (!tearSlot) return gear

  const tearSlotItem = gear.getEquippedItem(tearSlot);

  for (const [socketIdx, socketColor] of tearSlotItem!.allSocketColors().entries()) {
    if (socketColor === tearColor) {
      return gear.withEquippedItem(tearSlot, tearSlotItem!.withGem(sim.db.lookupGem(49110), socketIdx), true);
    }
  }

  return gear;
}

const findSocketsByColor = (player: Player<AutoGemSpec>, gear: Gear, epWeights: Stats, color: GemColor, tearSlot: ItemSlot | null): Array<[ItemSlot, number]> => {
  const socketList = new Array<[ItemSlot, number]>();
  const isBlacksmithing = player.isBlacksmithing();

  for (var slot of gear.getItemSlots()) {
    const item = gear.getEquippedItem(slot);

    if (!item) {
      continue;
    }

    const ignoreYellowSockets = ((item!.numSocketsOfColor(GemColor.GemColorBlue) > 0) && (slot !== tearSlot))

    for (const [socketIdx, socketColor] of item!.curSocketColors(isBlacksmithing).entries()) {
      if (item!.hasSocketedGem(socketIdx)) {
        continue;
      }

      let matchYellowSocket = false;

      if ((socketColor === GemColor.GemColorYellow) && !ignoreYellowSockets) {
        matchYellowSocket = new Stats(item.item.socketBonus).computeEP(epWeights) > 1e-8;
      }

      if (((color === GemColor.GemColorYellow) && matchYellowSocket) || ((color === GemColor.GemColorRed) && !matchYellowSocket)) {
        socketList.push([slot, socketIdx]);
      }
    }
  }

  return socketList;
}

/**
 * Determine if player is trying to reach hard cap
 * @remarks
 * Used for Feral sim only
 */
const detectArpStackConfiguration = (player: Player<Spec.SpecFeralDruid>, arpCap: number, arpTarget: number): boolean => {
  const currentArp = Stats.fromProto(player.getCurrentStats().finalStats).getStat(Stat.StatArmorPenetration);
  return (arpTarget > 1000) && (currentArp > 648) && ((currentArp + 20) < arpCap);
}

const sortYellowSockets = (gear: Gear, yellowSocketList: Array<[ItemSlot, number]>, epWeights: Stats, tearSlot: ItemSlot | null) => {
  return yellowSocketList.sort(([slot1], [slot2]) => {

    // If both yellow sockets belong to the same item, then treat them equally.
    if (slot1 === slot2) {
      return 0;
    }

    // If an item already has a Nightmare Tear socketed, then bump up any yellow sockets in it to highest priority.
    if (slot1 === tearSlot) {
      return -1;
    }

    if (slot2 === tearSlot) {
      return 1;
    }

    // For all other cases, sort by the ratio of the socket bonus value divided by the number of yellow sockets required to activate it.
    const item1 = gear.getEquippedItem(slot1);
    const bonus1 = new Stats(item1?.item.socketBonus).computeEP(epWeights);
    const item2 = gear.getEquippedItem(slot2);
    const bonus2 = new Stats(item2?.item.socketBonus).computeEP(epWeights);
    return bonus2 / (item2?.numSocketsOfColor(GemColor.GemColorYellow) || 0) - bonus1 / (item1?.numSocketsOfColor(GemColor.GemColorYellow) || 0);
  }).reverse();
}

export const optimizeGems = async (sim: Sim, player: Player<AutoGemSpec>) => {
  // First, clear all existing gems
  let optimizedGear = player.getGear().withoutGems();

  // Next, socket the meta
  switch (player.spec) {
    // Same for all specs
    case Spec.SpecFeralDruid:
    case Spec.SpecHunter:
    case Spec.SpecWarrior:
      optimizedGear = optimizedGear.withMetaGem(sim.db.lookupGem(41398));
  }

  // Next, socket a Nightmare Tear in the best socket
  const epWeights = player.getEpWeights();
  let tearColor = GemColor.GemColorBlue;
  let tearSlot = findBlueTearSlot(optimizedGear, epWeights);
  if (tearSlot === null) {
    tearColor = GemColor.GemColorYellow;
    tearSlot = findYellowTearSlot(optimizedGear, epWeights);
  }
  optimizedGear = socketTear(sim, optimizedGear, tearSlot, tearColor);
  await updateGear(sim, player, optimizedGear);

  let arpTarget = calcArpTarget(optimizedGear).getStat(Stat.StatArmorPenetration);
  const arpCap = calcArpCap(arpTarget, player);
  const expCap = calcExpCap(player);
  const critCap = calcCritCap(optimizedGear);
  const hitCap = calcHitCap(player);

  // Should we gem expertise?
  const enableExpertiseGemming = player.spec === Spec.SpecFeralDruid || (player.spec === Spec.SpecWarrior && !player.getDisableExpertiseGemming())

  // Next, identify all sockets where red gems will be placed
  const redSockets = findSocketsByColor(player, optimizedGear, epWeights, GemColor.GemColorRed, tearSlot);
  // Rank order red gems to use with their associated stat caps
  const redGemCaps = new Array<[number, Stats]>();
  redGemCaps.push([GemsByStats.Arp, arpCap]);

  if (enableExpertiseGemming) {
    redGemCaps.push([GemsByStats.Exp, expCap]);
  }

  // If Feral swap to Agi if below crit cap
  if (player.spec === Spec.SpecFeralDruid) {
    redGemCaps.push([GemsByStats.Agi, critCap]);
  }
  // If Feral or Warrior swap to Str when ArP and Crit capped
  if (player.spec === Spec.SpecFeralDruid || player.spec === Spec.SpecWarrior) {
    redGemCaps.push([GemsByStats.Str, new Stats()]);
  }
  // If Hunter swap to Agi when ArP and Crit capped
  if (player.spec === Spec.SpecHunter) {
    redGemCaps.push([GemsByStats.Agi, new Stats()]);
  }

  // If JC, then socket 34 ArP gems in first three red sockets before proceeding
  let startIdx = 0;
  if (player.hasProfession(Profession.Jewelcrafting)) {
    optimizedGear = optimizeJewelCraftingGems(sim, player, optimizedGear, redSockets, arpCap, arpTarget);
    startIdx = 3;
  }

  // Do multiple passes to fill in red gems up their caps
  optimizedGear = await fillGemsToCaps(sim, player, optimizedGear, redSockets, redGemCaps, 0, startIdx);

  // Now repeat the process for yellow gems
  const yellowSockets = findSocketsByColor(player, optimizedGear, epWeights, GemColor.GemColorYellow, tearSlot);
  const yellowGemCaps = new Array<[number, Stats]>();

  // Rigid Ametrine
  yellowGemCaps.push([GemsByStats.Hit, hitCap]);

  switch (player.spec) {
    case Spec.SpecFeralDruid:
      yellowGemCaps.push([GemsByStats.Exp_Hit, hitCap.add(expCap)]);

      // Allow for socketing ArP gems in weaker yellow sockets after capping Hit and Expertise
      // when ArP stacking is detected
      if (detectArpStackConfiguration(player, arpCap.getStat(Stat.StatArmorPenetration), arpTarget)) {
        sortYellowSockets(optimizedGear, yellowSockets, epWeights, tearSlot);
        yellowGemCaps.push([GemsByStats.Arp, arpCap]);
      }

      yellowGemCaps.push([GemsByStats.Agi_Hit, hitCap.add(critCap)]);
      yellowGemCaps.push([GemsByStats.Str_Hit, hitCap]);
      yellowGemCaps.push([GemsByStats.Agi_Crit, critCap]);
      yellowGemCaps.push([GemsByStats.Str_Crit, critCap]);
      yellowGemCaps.push([GemsByStats.Str_Haste, new Stats()]);
      break
    case Spec.SpecWarrior:
      if (enableExpertiseGemming) yellowGemCaps.push([GemsByStats.Exp_Hit, hitCap.add(expCap)]);

      // Allow for socketing ArP gems in weaker yellow sockets after capping Hit and Expertise
      // when ArP stacking is detected
      if (detectArpStackConfiguration(player, arpCap.getStat(Stat.StatArmorPenetration), arpTarget)) {
        sortYellowSockets(optimizedGear, yellowSockets, epWeights, tearSlot);
        yellowGemCaps.push([GemsByStats.Arp, arpCap]);
      }

      yellowGemCaps.push([GemsByStats.Str_Hit, hitCap]);
      yellowGemCaps.push([GemsByStats.Str_Crit, critCap]);
      break
    case Spec.SpecHunter:
      yellowGemCaps.push([GemsByStats.Agi_Hit, hitCap]);
      yellowGemCaps.push([GemsByStats.Agi_Crit, critCap]);
      break
  }


  await fillGemsToCaps(sim, player, optimizedGear, yellowSockets, yellowGemCaps, 0, 0);
}

