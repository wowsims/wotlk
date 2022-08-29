package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"
)

type Stats [42]float64

type ItemResponse interface {
	GetName() string
	GetQuality() int
	GetIcon() string
	TooltipWithoutSetBonus() string
	GetTooltipRegexString(pattern *regexp.Regexp, matchIdx int) string
	GetTooltipRegexValue(pattern *regexp.Regexp, matchIdx int) int
	GetIntValue(pattern *regexp.Regexp) int
	GetStats() Stats
	GetClassAllowlist() []proto.Class
	IsEquippable() bool
	GetItemLevel() int
	GetPhase() int
	GetUnique() bool
	GetItemType() proto.ItemType
	GetArmorType() proto.ArmorType
	GetWeaponType() proto.WeaponType
	GetHandType() proto.HandType
	GetRangedWeaponType() proto.RangedWeaponType
	GetWeaponDamage() (float64, float64)
	GetWeaponSpeed() float64
	GetGemSockets() []proto.GemColor
	GetSocketBonus() Stats
	GetSocketColor() proto.GemColor
	GetGemStats() Stats
	GetItemSetName() string
	IsHeroic() bool
	GetRequiredProfession() proto.Profession
}

type WowheadItemResponse struct {
	Name    string `json:"name"`
	Quality int    `json:"quality"`
	Icon    string `json:"icon"`
	Tooltip string `json:"tooltip"`
}

func (item WowheadItemResponse) GetName() string {
	return item.Name
}
func (item WowheadItemResponse) GetQuality() int {
	return item.Quality
}
func (item WowheadItemResponse) GetIcon() string {
	return item.Icon
}

func GetRegexStringValue(srcStr string, pattern *regexp.Regexp, matchIdx int) string {
	match := pattern.FindStringSubmatch(srcStr)
	if match == nil {
		return ""
	} else {
		return match[matchIdx]
	}
}
func GetRegexIntValue(srcStr string, pattern *regexp.Regexp, matchIdx int) int {
	matchStr := GetRegexStringValue(srcStr, pattern, matchIdx)

	val, err := strconv.Atoi(matchStr)
	if err != nil {
		return 0
	}

	return val
}
func GetBestRegexIntValue(srcStr string, patterns []*regexp.Regexp, matchIdx int) int {
	best := 0
	for _, pattern := range patterns {
		newVal := GetRegexIntValue(srcStr, pattern, matchIdx)
		if newVal > best {
			best = newVal
		}
	}
	return best
}

func (item WowheadItemResponse) TooltipWithoutSetBonus() string {
	setIdx := strings.Index(item.Tooltip, "Set : ")
	if setIdx == -1 {
		return item.Tooltip
	} else {
		return item.Tooltip[:setIdx]
	}
}

func (item WowheadItemResponse) GetTooltipRegexString(pattern *regexp.Regexp, matchIdx int) string {
	return GetRegexStringValue(item.TooltipWithoutSetBonus(), pattern, matchIdx)
}

func (item WowheadItemResponse) GetTooltipRegexValue(pattern *regexp.Regexp, matchIdx int) int {
	return GetRegexIntValue(item.TooltipWithoutSetBonus(), pattern, matchIdx)
}

func (item WowheadItemResponse) GetIntValue(pattern *regexp.Regexp) int {
	return item.GetTooltipRegexValue(pattern, 1)
}

var armorRegex = regexp.MustCompile("<!--amr-->([0-9]+) Armor")
var agilityRegex = regexp.MustCompile("<!--stat3-->\\+([0-9]+) Agility")
var strengthRegex = regexp.MustCompile("<!--stat4-->\\+([0-9]+) Strength")
var intellectRegex = regexp.MustCompile("<!--stat5-->\\+([0-9]+) Intellect")
var spiritRegex = regexp.MustCompile("<!--stat6-->\\+([0-9]+) Spirit")
var staminaRegex = regexp.MustCompile("<!--stat7-->\\+([0-9]+) Stamina")
var spellPowerRegex = regexp.MustCompile("Increases spell power by ([0-9]+)\\.")
var spellPowerRegex2 = regexp.MustCompile("Increases spell power by <!--rtg45-->([0-9]+)\\.")

// Not sure these exist anymore?
var arcaneSpellPowerRegex = regexp.MustCompile("Increases Arcane power by ([0-9]+)\\.")
var fireSpellPowerRegex = regexp.MustCompile("Increases Fire power by ([0-9]+)\\.")
var frostSpellPowerRegex = regexp.MustCompile("Increases Frost power by ([0-9]+)\\.")
var holySpellPowerRegex = regexp.MustCompile("Increases Holy power by ([0-9]+)\\.")
var natureSpellPowerRegex = regexp.MustCompile("Increases Nature power by ([0-9]+)\\.")
var shadowSpellPowerRegex = regexp.MustCompile("Increases Shadow power by ([0-9]+)\\.")

var hitRegex = regexp.MustCompile("Improves hit rating by <!--rtg31-->([0-9]+)\\.")
var critRegex = regexp.MustCompile("Improves critical strike rating by <!--rtg32-->([0-9]+)\\.")
var hasteRegex = regexp.MustCompile("Improves haste rating by <!--rtg36-->([0-9]+)\\.")

var spellPenetrationRegex = regexp.MustCompile("Increases your spell penetration by ([0-9]+)\\.")
var mp5Regex = regexp.MustCompile("Restores ([0-9]+) mana per 5 sec\\.")
var attackPowerRegex = regexp.MustCompile("Increases attack power by ([0-9]+)\\.")
var attackPowerRegex2 = regexp.MustCompile("Increases attack power by <!--rtg38-->([0-9]+)\\.")

var rangedAttackPowerRegex = regexp.MustCompile("Increases ranged attack power by ([0-9]+)\\.")
var rangedAttackPowerRegex2 = regexp.MustCompile("Increases ranged attack power by <!--rtg39-->([0-9]+)\\.")

var armorPenetrationRegex = regexp.MustCompile("Increases armor penetration rating by ([0-9]+)")
var armorPenetrationRegex2 = regexp.MustCompile("Increases your armor penetration by <!--rtg44-->([0-9]+)\\.")

var expertiseRegex = regexp.MustCompile("Increases your expertise rating by <!--rtg37-->([0-9]+)\\.")
var weaponDamageRegex = regexp.MustCompile("<!--dmg-->([0-9]+) - ([0-9]+)")
var weaponDamageRegex2 = regexp.MustCompile("<!--dmg-->([0-9]+) Damage")
var weaponSpeedRegex = regexp.MustCompile("<!--spd-->(([0-9]+).([0-9]+))")

var defenseRegex = regexp.MustCompile("Increases defense rating by <!--rtg12-->([0-9]+)\\.")
var defenseRegex2 = regexp.MustCompile("Increases defense rating by ([0-9]+)\\.")
var blockRegex = regexp.MustCompile("Increases your shield block rating by <!--rtg15-->([0-9]+)\\.")
var blockRegex2 = regexp.MustCompile("Increases your shield block rating by ([0-9]+)\\.")
var blockValueRegex = regexp.MustCompile("Increases the block value of your shield by ([0-9]+)\\.")
var blockValueRegex2 = regexp.MustCompile("<br>([0-9]+) Block<br>")
var dodgeRegex = regexp.MustCompile("Increases your dodge rating by <!--rtg13-->([0-9]+)\\.")
var dodgeRegex2 = regexp.MustCompile("Increases your dodge rating by ([0-9]+)\\.")
var parryRegex = regexp.MustCompile("Increases your parry rating by <!--rtg14-->([0-9]+)\\.")
var parryRegex2 = regexp.MustCompile("Increases your parry rating by ([0-9]+)\\.")
var resilienceRegex = regexp.MustCompile("Improves your resilience rating by <!--rtg35-->([0-9]+)\\.")
var arcaneResistanceRegex = regexp.MustCompile("\\+([0-9]+) Arcane Resistance")
var fireResistanceRegex = regexp.MustCompile("\\+([0-9]+) Fire Resistance")
var frostResistanceRegex = regexp.MustCompile("\\+([0-9]+) Frost Resistance")
var natureResistanceRegex = regexp.MustCompile("\\+([0-9]+) Nature Resistance")
var shadowResistanceRegex = regexp.MustCompile("\\+([0-9]+) Shadow Resistance")

func (item WowheadItemResponse) GetStats() Stats {
	sp := float64(item.GetIntValue(spellPowerRegex)) + float64(item.GetIntValue(spellPowerRegex2))
	baseAP := float64(item.GetIntValue(attackPowerRegex)) + float64(item.GetIntValue(attackPowerRegex2))
	return Stats{
		proto.Stat_StatArmor:             float64(item.GetIntValue(armorRegex)),
		proto.Stat_StatStrength:          float64(item.GetIntValue(strengthRegex)),
		proto.Stat_StatAgility:           float64(item.GetIntValue(agilityRegex)),
		proto.Stat_StatStamina:           float64(item.GetIntValue(staminaRegex)),
		proto.Stat_StatIntellect:         float64(item.GetIntValue(intellectRegex)),
		proto.Stat_StatSpirit:            float64(item.GetIntValue(spiritRegex)),
		proto.Stat_StatSpellPower:        sp,
		proto.Stat_StatHealingPower:      sp,
		proto.Stat_StatArcaneSpellPower:  float64(item.GetIntValue(arcaneSpellPowerRegex)),
		proto.Stat_StatFireSpellPower:    float64(item.GetIntValue(fireSpellPowerRegex)),
		proto.Stat_StatFrostSpellPower:   float64(item.GetIntValue(frostSpellPowerRegex)),
		proto.Stat_StatHolySpellPower:    float64(item.GetIntValue(holySpellPowerRegex)),
		proto.Stat_StatNatureSpellPower:  float64(item.GetIntValue(natureSpellPowerRegex)),
		proto.Stat_StatShadowSpellPower:  float64(item.GetIntValue(shadowSpellPowerRegex)),
		proto.Stat_StatSpellHit:          float64(item.GetIntValue(hitRegex)),
		proto.Stat_StatMeleeHit:          float64(item.GetIntValue(hitRegex)),
		proto.Stat_StatSpellCrit:         float64(item.GetIntValue(critRegex)),
		proto.Stat_StatMeleeCrit:         float64(item.GetIntValue(critRegex)),
		proto.Stat_StatSpellHaste:        float64(item.GetIntValue(hasteRegex)),
		proto.Stat_StatMeleeHaste:        float64(item.GetIntValue(hasteRegex)),
		proto.Stat_StatSpellPenetration:  float64(item.GetIntValue(spellPenetrationRegex)),
		proto.Stat_StatMP5:               float64(item.GetIntValue(mp5Regex)),
		proto.Stat_StatAttackPower:       baseAP,
		proto.Stat_StatRangedAttackPower: baseAP + float64(item.GetIntValue(rangedAttackPowerRegex)) + float64(item.GetIntValue(rangedAttackPowerRegex2)),
		proto.Stat_StatArmorPenetration:  float64(item.GetIntValue(armorPenetrationRegex) + item.GetIntValue(armorPenetrationRegex2)),
		proto.Stat_StatExpertise:         float64(item.GetIntValue(expertiseRegex)),
		proto.Stat_StatDefense:           float64(item.GetIntValue(defenseRegex) + item.GetIntValue(defenseRegex2)),
		proto.Stat_StatBlock:             float64(item.GetIntValue(blockRegex) + item.GetIntValue(blockRegex2)),
		proto.Stat_StatBlockValue:        float64(item.GetIntValue(blockValueRegex) + item.GetIntValue(blockValueRegex2)),
		proto.Stat_StatDodge:             float64(item.GetIntValue(dodgeRegex) + item.GetIntValue(dodgeRegex2)),
		proto.Stat_StatParry:             float64(item.GetIntValue(parryRegex) + item.GetIntValue(parryRegex2)),
		proto.Stat_StatResilience:        float64(item.GetIntValue(resilienceRegex)),
		proto.Stat_StatArcaneResistance:  float64(item.GetIntValue(arcaneResistanceRegex)),
		proto.Stat_StatFireResistance:    float64(item.GetIntValue(fireResistanceRegex)),
		proto.Stat_StatFrostResistance:   float64(item.GetIntValue(frostResistanceRegex)),
		proto.Stat_StatNatureResistance:  float64(item.GetIntValue(natureResistanceRegex)),
		proto.Stat_StatShadowResistance:  float64(item.GetIntValue(shadowResistanceRegex)),
	}
}

type classPattern struct {
	class   proto.Class
	pattern *regexp.Regexp
}

// Detects class-locked items, e.g. tier sets and pvp gear.
var classPatternsWowhead = []classPattern{
	{class: proto.Class_ClassWarrior, pattern: regexp.MustCompile(`<a href=\"/wotlk/warrior\" class=\"c1\">Warrior</a>`)},
	{class: proto.Class_ClassPaladin, pattern: regexp.MustCompile(`<a href=\"/wotlk/paladin\" class=\"c2\">Paladin</a>`)},
	{class: proto.Class_ClassHunter, pattern: regexp.MustCompile(`<a href=\"/wotlk/hunter\" class=\"c3\">Hunter</a>`)},
	{class: proto.Class_ClassRogue, pattern: regexp.MustCompile(`<a href=\"/wotlk/rogue\" class=\"c4\">Rogue</a>`)},
	{class: proto.Class_ClassPriest, pattern: regexp.MustCompile(`<a href=\"/wotlk/priest\" class=\"c5\">Priest</a>`)},
	{class: proto.Class_ClassDeathknight, pattern: regexp.MustCompile(`<a href=\"/wotlk/death-knight\" class=\"c6\">Death Knight</a>`)},
	{class: proto.Class_ClassShaman, pattern: regexp.MustCompile(`<a href=\"/wotlk/shaman\" class=\"c7\">Shaman</a>`)},
	{class: proto.Class_ClassMage, pattern: regexp.MustCompile(`<a href=\"/wotlk/mage\" class=\"c8\">Mage</a>`)},
	{class: proto.Class_ClassWarlock, pattern: regexp.MustCompile(`<a href=\"/wotlk/warlock\" class=\"c9\">Warlock</a>`)},
	{class: proto.Class_ClassDruid, pattern: regexp.MustCompile(`<a href=\"/wotlk/druid\" class=\"c11\">Druid</a>`)},
}

func (item WowheadItemResponse) GetClassAllowlist() []proto.Class {
	var allowlist []proto.Class

	for _, entry := range classPatternsWowhead {
		if entry.pattern.MatchString(item.Tooltip) {
			allowlist = append(allowlist, entry.class)
		}
	}

	return allowlist
}

// At least one of these regexes must be present for the item to be equippable.
var requiredEquippableRegexes = []*regexp.Regexp{
	regexp.MustCompile("<td>Head</td>"),
	regexp.MustCompile("<td>Neck</td>"),
	regexp.MustCompile("<td>Shoulder</td>"),
	regexp.MustCompile("<td>Back</td>"),
	regexp.MustCompile("<td>Chest</td>"),
	regexp.MustCompile("<td>Wrist</td>"),
	regexp.MustCompile("<td>Hands</td>"),
	regexp.MustCompile("<td>Waist</td>"),
	regexp.MustCompile("<td>Legs</td>"),
	regexp.MustCompile("<td>Feet</td>"),
	regexp.MustCompile("<td>Finger</td>"),
	regexp.MustCompile("<td>Trinket</td>"),
	regexp.MustCompile("<td>Ranged</td>"),
	regexp.MustCompile("<td>Thrown</td>"),
	regexp.MustCompile("<td>Relic</td>"),
	regexp.MustCompile("<td>Main Hand</td>"),
	regexp.MustCompile("<td>Two-Hand</td>"),
	regexp.MustCompile("<td>One-Hand</td>"),
	regexp.MustCompile("<td>Off Hand</td>"),
	regexp.MustCompile("<td>Held In Off-hand</td>"),
	regexp.MustCompile("<td>Held In Off-Hand</td>"),
}

// If any of these regexes are present, the item is not equippable.
var nonEquippableRegexes = []*regexp.Regexp{
	regexp.MustCompile("Design:"),
	regexp.MustCompile("Recipe:"),
	regexp.MustCompile("Pattern:"),
	regexp.MustCompile("Plans:"),
	regexp.MustCompile("Schematic:"),
	regexp.MustCompile("Random enchantment"),
}

func (item WowheadItemResponse) IsEquippable() bool {
	found := false
	for _, pattern := range requiredEquippableRegexes {
		if pattern.MatchString(item.Tooltip) {
			found = true
		}
	}
	if !found {
		return false
	}

	for _, pattern := range nonEquippableRegexes {
		if pattern.MatchString(item.Tooltip) {
			return false
		}
	}

	return true
}

func (item WowheadItemResponse) IsPattern() bool {
	for _, pattern := range nonEquippableRegexes {
		if pattern.MatchString(item.Tooltip) {
			return true
		}
	}
	return false
}

var itemLevelRegex = regexp.MustCompile("Item Level <!--ilvl-->([0-9]+)<")

func (item WowheadItemResponse) GetItemLevel() int {
	return item.GetIntValue(itemLevelRegex)
}

var phaseRegex = regexp.MustCompile("Phase ([0-9])")

func (item WowheadItemResponse) GetPhase() int {
	phase := item.GetIntValue(phaseRegex)
	if phase != 0 {
		return phase
	}

	ilvl := item.GetItemLevel()
	if ilvl <= 164 { // TBC items
		return 0
	}

	if ilvl < 200 || ilvl == 200 || ilvl == 213 || ilvl == 226 {
		return 1
	} else if ilvl == 219 || ilvl == 226 || ilvl == 239 {
		return 2
	} else if ilvl == 232 || ilvl == 245 || ilvl == 258 {
		return 3
	} else if ilvl == 251 || ilvl == 258 || ilvl == 259 || ilvl == 264 || ilvl == 268 || ilvl == 270 || ilvl == 271 || ilvl == 272 {
		return 4
	} else if ilvl == 277 || ilvl == 284 {
		return 5
	}

	// default to 1
	return 1
}

var uniqueRegex = regexp.MustCompile("Unique")
var jcGemsRegex = regexp.MustCompile("Jeweler's Gems")

func (item WowheadItemResponse) GetUnique() bool {
	return uniqueRegex.MatchString(item.Tooltip) && !jcGemsRegex.MatchString(item.Tooltip)
}

var itemTypePatterns = map[proto.ItemType]*regexp.Regexp{
	proto.ItemType_ItemTypeHead:     regexp.MustCompile("<td>Head</td>"),
	proto.ItemType_ItemTypeNeck:     regexp.MustCompile("<td>Neck</td>"),
	proto.ItemType_ItemTypeShoulder: regexp.MustCompile("<td>Shoulder</td>"),
	proto.ItemType_ItemTypeBack:     regexp.MustCompile("<td>Back</td>"),
	proto.ItemType_ItemTypeChest:    regexp.MustCompile("<td>Chest</td>"),
	proto.ItemType_ItemTypeWrist:    regexp.MustCompile("<td>Wrist</td>"),
	proto.ItemType_ItemTypeHands:    regexp.MustCompile("<td>Hands</td>"),
	proto.ItemType_ItemTypeWaist:    regexp.MustCompile("<td>Waist</td>"),
	proto.ItemType_ItemTypeLegs:     regexp.MustCompile("<td>Legs</td>"),
	proto.ItemType_ItemTypeFeet:     regexp.MustCompile("<td>Feet</td>"),
	proto.ItemType_ItemTypeFinger:   regexp.MustCompile("<td>Finger</td>"),
	proto.ItemType_ItemTypeTrinket:  regexp.MustCompile("<td>Trinket</td>"),
	proto.ItemType_ItemTypeWeapon:   regexp.MustCompile("<td>((Main Hand)|(Two-Hand)|(One-Hand)|(Off Hand)|(Held In Off-hand)|(Held In Off-Hand))</td>"),
	proto.ItemType_ItemTypeRanged:   regexp.MustCompile("<td>(Ranged|Thrown|Relic)</td>"),
}

func (item WowheadItemResponse) GetItemType() proto.ItemType {
	for itemType, pattern := range itemTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return itemType
		}
	}
	panic("Could not find item type from tooltip: " + item.Tooltip)
}

var armorTypePatterns = map[proto.ArmorType]*regexp.Regexp{
	proto.ArmorType_ArmorTypeCloth:   regexp.MustCompile("<span class=\\\"q1\\\">Cloth</span>"),
	proto.ArmorType_ArmorTypeLeather: regexp.MustCompile("<span class=\\\"q1\\\">Leather</span>"),
	proto.ArmorType_ArmorTypeMail:    regexp.MustCompile("<span class=\\\"q1\\\">Mail</span>"),
	proto.ArmorType_ArmorTypePlate:   regexp.MustCompile("<span class=\\\"q1\\\">Plate</span>"),
}

func (item WowheadItemResponse) GetArmorType() proto.ArmorType {
	for armorType, pattern := range armorTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return armorType
		}
	}
	return proto.ArmorType_ArmorTypeUnknown
}

var weaponTypePatterns = map[proto.WeaponType]*regexp.Regexp{
	proto.WeaponType_WeaponTypeAxe:     regexp.MustCompile("<span class=\\\"q1\\\">Axe</span>"),
	proto.WeaponType_WeaponTypeDagger:  regexp.MustCompile("<span class=\\\"q1\\\">Dagger</span>"),
	proto.WeaponType_WeaponTypeFist:    regexp.MustCompile("<span class=\\\"q1\\\">Fist Weapon</span>"),
	proto.WeaponType_WeaponTypeMace:    regexp.MustCompile("<span class=\\\"q1\\\">Mace</span>"),
	proto.WeaponType_WeaponTypeOffHand: regexp.MustCompile("<td>Held In Off-hand</td>"),
	proto.WeaponType_WeaponTypePolearm: regexp.MustCompile("<span class=\\\"q1\\\">Polearm</span>"),
	proto.WeaponType_WeaponTypeShield:  regexp.MustCompile("<span class=\\\"q1\\\">Shield</span>"),
	proto.WeaponType_WeaponTypeStaff:   regexp.MustCompile("<span class=\\\"q1\\\">Staff</span>"),
	proto.WeaponType_WeaponTypeSword:   regexp.MustCompile("<span class=\\\"q1\\\">Sword</span>"),
}

func (item WowheadItemResponse) GetWeaponType() proto.WeaponType {
	for weaponType, pattern := range weaponTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return weaponType
		}
	}
	return proto.WeaponType_WeaponTypeUnknown
}

var handTypePatterns = map[proto.HandType]*regexp.Regexp{
	proto.HandType_HandTypeMainHand: regexp.MustCompile("<td>Main Hand</td>"),
	proto.HandType_HandTypeOneHand:  regexp.MustCompile("<td>One-Hand</td>"),
	proto.HandType_HandTypeOffHand:  regexp.MustCompile("<td>((Off Hand)|(Held In Off-hand)|(Held In Off-Hand))</td>"),
	proto.HandType_HandTypeTwoHand:  regexp.MustCompile("<td>Two-Hand</td>"),
}

func (item WowheadItemResponse) GetHandType() proto.HandType {
	for handType, pattern := range handTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return handType
		}
	}
	return proto.HandType_HandTypeUnknown
}

var rangedWeaponTypePatterns = map[proto.RangedWeaponType]*regexp.Regexp{
	proto.RangedWeaponType_RangedWeaponTypeBow:      regexp.MustCompile("<span class=\\\"q1\\\">Bow</span>"),
	proto.RangedWeaponType_RangedWeaponTypeCrossbow: regexp.MustCompile("<span class=\\\"q1\\\">Crossbow</span>"),
	proto.RangedWeaponType_RangedWeaponTypeGun:      regexp.MustCompile("<span class=\\\"q1\\\">Gun</span>"),
	proto.RangedWeaponType_RangedWeaponTypeIdol:     regexp.MustCompile("<span class=\\\"q1\\\">Idol</span>"),
	proto.RangedWeaponType_RangedWeaponTypeLibram:   regexp.MustCompile("<span class=\\\"q1\\\">Libram</span>"),
	proto.RangedWeaponType_RangedWeaponTypeThrown:   regexp.MustCompile("<span class=\\\"q1\\\">Thrown</span>"),
	proto.RangedWeaponType_RangedWeaponTypeTotem:    regexp.MustCompile("<span class=\\\"q1\\\">Totem</span>"),
	proto.RangedWeaponType_RangedWeaponTypeWand:     regexp.MustCompile("<span class=\\\"q1\\\">Wand</span>"),
	proto.RangedWeaponType_RangedWeaponTypeSigil:    regexp.MustCompile("<span class=\\\"q1\\\">Sigil</span>"),
}

func (item WowheadItemResponse) GetRangedWeaponType() proto.RangedWeaponType {
	for rangedWeaponType, pattern := range rangedWeaponTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return rangedWeaponType
		}
	}
	return proto.RangedWeaponType_RangedWeaponTypeUnknown
}

// Returns min/max of weapon damage
func (item WowheadItemResponse) GetWeaponDamage() (float64, float64) {
	noCommas := strings.ReplaceAll(item.Tooltip, ",", "")
	if matches := weaponDamageRegex.FindStringSubmatch(noCommas); len(matches) > 0 {
		min, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		max, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		if min > max {
			log.Fatalf("Invalid weapon damage for item %s: min = %0.1f, max = %0.1f", item.Name, min, max)
		}
		return min, max
	} else if matches := weaponDamageRegex2.FindStringSubmatch(noCommas); len(matches) > 0 {
		val, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		return val, val
	}
	return 0, 0
}

func (item WowheadItemResponse) GetWeaponSpeed() float64 {
	if matches := weaponSpeedRegex.FindStringSubmatch(item.Tooltip); len(matches) > 0 {
		speed, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		return speed
	}
	return 0
}

var gemColorsRegex, _ = regexp.Compile("(Meta|Yellow|Blue|Red) Socket")

func (item WowheadItemResponse) GetGemSockets() []proto.GemColor {
	matches := gemColorsRegex.FindAllStringSubmatch(item.Tooltip, -1)
	if matches == nil {
		return []proto.GemColor{}
	}

	numSockets := len(matches)
	gemColors := make([]proto.GemColor, numSockets)
	for socketIdx, match := range matches {
		gemColorName := "GemColor" + match[1]
		gemColors[socketIdx] = proto.GemColor(proto.GemColor_value[gemColorName])
	}
	return gemColors
}

var socketBonusRegex = regexp.MustCompile("<span class=\\\"q0\\\">Socket Bonus: (.*?)</span>")
var strengthSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Strength")}
var agilitySocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Agility")}
var staminaSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Stamina")}
var intellectSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Intellect")}
var spiritSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Spirit")}
var spellPowerSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Spell Power")}
var spellHitSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Hit Rating")}
var spellCritSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Critical Strike Rating")}
var hasteSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Haste Rating")}
var mp5SocketBonusRegexes = []*regexp.Regexp{
	regexp.MustCompile("([0-9]+) Mana per 5 sec"),
	regexp.MustCompile("([0-9]+) mana per 5 sec"),
}
var attackPowerSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Attack Power")}
var armorPenSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Armor Penetration Rating")}
var expertiseSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Expertise Rating")}
var defenseSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Defense Rating")}
var blockSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Block Rating")}
var dodgeSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Dodge Rating")}
var parrySocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Parry Rating")}
var resilienceSocketBonusRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Resilience Rating")}

func (item WowheadItemResponse) GetSocketBonus() Stats {
	match := socketBonusRegex.FindStringSubmatch(item.Tooltip)
	if match == nil {
		return Stats{}
	}

	bonusStr := match[1]
	//fmt.Printf("\n%s\n", bonusStr)

	stats := Stats{
		proto.Stat_StatStrength:          float64(GetBestRegexIntValue(bonusStr, strengthSocketBonusRegexes, 1)),
		proto.Stat_StatAgility:           float64(GetBestRegexIntValue(bonusStr, agilitySocketBonusRegexes, 1)),
		proto.Stat_StatStamina:           float64(GetBestRegexIntValue(bonusStr, staminaSocketBonusRegexes, 1)),
		proto.Stat_StatIntellect:         float64(GetBestRegexIntValue(bonusStr, intellectSocketBonusRegexes, 1)),
		proto.Stat_StatSpirit:            float64(GetBestRegexIntValue(bonusStr, spiritSocketBonusRegexes, 1)),
		proto.Stat_StatSpellHaste:        float64(GetBestRegexIntValue(bonusStr, hasteSocketBonusRegexes, 1)),
		proto.Stat_StatSpellPower:        float64(GetBestRegexIntValue(bonusStr, spellPowerSocketBonusRegexes, 1)),
		proto.Stat_StatHealingPower:      float64(GetBestRegexIntValue(bonusStr, spellPowerSocketBonusRegexes, 1)),
		proto.Stat_StatSpellHit:          float64(GetBestRegexIntValue(bonusStr, spellHitSocketBonusRegexes, 1)),
		proto.Stat_StatMeleeHit:          float64(GetBestRegexIntValue(bonusStr, spellHitSocketBonusRegexes, 1)),
		proto.Stat_StatSpellCrit:         float64(GetBestRegexIntValue(bonusStr, spellCritSocketBonusRegexes, 1)),
		proto.Stat_StatMeleeCrit:         float64(GetBestRegexIntValue(bonusStr, spellCritSocketBonusRegexes, 1)),
		proto.Stat_StatMeleeHaste:        float64(GetBestRegexIntValue(bonusStr, hasteSocketBonusRegexes, 1)),
		proto.Stat_StatMP5:               float64(GetBestRegexIntValue(bonusStr, mp5SocketBonusRegexes, 1)),
		proto.Stat_StatAttackPower:       float64(GetBestRegexIntValue(bonusStr, attackPowerSocketBonusRegexes, 1)),
		proto.Stat_StatRangedAttackPower: float64(GetBestRegexIntValue(bonusStr, attackPowerSocketBonusRegexes, 1)),
		proto.Stat_StatExpertise:         float64(GetBestRegexIntValue(bonusStr, expertiseSocketBonusRegexes, 1)),
		proto.Stat_StatArmorPenetration:  float64(GetBestRegexIntValue(bonusStr, armorPenSocketBonusRegexes, 1)),
		proto.Stat_StatDefense:           float64(GetBestRegexIntValue(bonusStr, defenseSocketBonusRegexes, 1)),
		proto.Stat_StatBlock:             float64(GetBestRegexIntValue(bonusStr, blockSocketBonusRegexes, 1)),
		proto.Stat_StatDodge:             float64(GetBestRegexIntValue(bonusStr, dodgeSocketBonusRegexes, 1)),
		proto.Stat_StatParry:             float64(GetBestRegexIntValue(bonusStr, parrySocketBonusRegexes, 1)),
		proto.Stat_StatResilience:        float64(GetBestRegexIntValue(bonusStr, resilienceSocketBonusRegexes, 1)),
	}

	return stats
}

var gemSocketColorPatterns = map[proto.GemColor]*regexp.Regexp{
	proto.GemColor_GemColorMeta:      regexp.MustCompile("Only fits in a meta gem slot\\."),
	proto.GemColor_GemColorBlue:      regexp.MustCompile("Matches a Blue (S|s)ocket\\."),
	proto.GemColor_GemColorRed:       regexp.MustCompile("Matches a Red (S|s)ocket\\."),
	proto.GemColor_GemColorYellow:    regexp.MustCompile("Matches a Yellow (S|s)ocket\\."),
	proto.GemColor_GemColorOrange:    regexp.MustCompile("Matches a ((Yellow)|(Red)) or ((Yellow)|(Red)) (S|s)ocket\\."),
	proto.GemColor_GemColorPurple:    regexp.MustCompile("Matches a ((Blue)|(Red)) or ((Blue)|(Red)) (S|s)ocket\\."),
	proto.GemColor_GemColorGreen:     regexp.MustCompile("Matches a ((Yellow)|(Blue)) or ((Yellow)|(Blue)) (S|s)ocket\\."),
	proto.GemColor_GemColorPrismatic: regexp.MustCompile("(Matches any (S|s)ocket)|(Matches a Red, Yellow or Blue (S|s)ocket)"),
}

func (item WowheadItemResponse) GetSocketColor() proto.GemColor {
	for socketColor, pattern := range gemSocketColorPatterns {
		if pattern.MatchString(item.Tooltip) {
			return socketColor
		}
	}
	// fmt.Printf("Could not find socket color for gem %s\n", item.Name)
	return proto.GemColor_GemColorUnknown
}

var strengthGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Strength"), regexp.MustCompile("\\+([0-9]+) (to )?All Stats")}
var agilityGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Agility"), regexp.MustCompile("\\+([0-9]+) (to )?All Stats")}
var staminaGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Stamina"), regexp.MustCompile("\\+([0-9]+) (to )?All Stats")}
var intellectGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Intellect"), regexp.MustCompile("\\+([0-9]+) (to )?All Stats")}
var spiritGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Spirit"), regexp.MustCompile("\\+([0-9]+) (to )?All Stats")}
var spellPowerGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Spell Power")}
var hitGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Hit Rating")}
var critGemStatRegexes = []*regexp.Regexp{
	regexp.MustCompile("\\+([0-9]+) Crit Rating"),
	regexp.MustCompile("\\+([0-9]+) Critical Strike Rating"),
	regexp.MustCompile("\\+([0-9]+) Critical"),
}
var hasteGemStatRegexes = []*regexp.Regexp{
	regexp.MustCompile("\\+([0-9]+) Haste Rating"),
}
var armorPenetrationGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Armor Penetration")}
var spellPenetrationGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Spell Penetration")}
var mp5GemStatRegexes = []*regexp.Regexp{
	regexp.MustCompile("([0-9]+) Mana per 5 sec"),
	regexp.MustCompile("([0-9]+) mana per 5 sec"),
	regexp.MustCompile("([0-9]+) Mana every 5 seconds"),
}
var attackPowerGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Attack Power")}
var expertiseGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Expertise Rating")}
var defenseGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Defense Rating")}
var dodgeGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Dodge Rating")}
var parryGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Parry Rating")}
var resilienceGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Resilience Rating")}
var allResistGemStatRegexes = []*regexp.Regexp{regexp.MustCompile("\\+([0-9]+) Resist All")}

func (item WowheadItemResponse) GetGemStats() Stats {
	stats := Stats{
		proto.Stat_StatStrength:  float64(GetBestRegexIntValue(item.Tooltip, strengthGemStatRegexes, 1)),
		proto.Stat_StatAgility:   float64(GetBestRegexIntValue(item.Tooltip, agilityGemStatRegexes, 1)),
		proto.Stat_StatStamina:   float64(GetBestRegexIntValue(item.Tooltip, staminaGemStatRegexes, 1)),
		proto.Stat_StatIntellect: float64(GetBestRegexIntValue(item.Tooltip, intellectGemStatRegexes, 1)),
		proto.Stat_StatSpirit:    float64(GetBestRegexIntValue(item.Tooltip, spiritGemStatRegexes, 1)),

		proto.Stat_StatSpellHit:   float64(GetBestRegexIntValue(item.Tooltip, hitGemStatRegexes, 1)),
		proto.Stat_StatMeleeHit:   float64(GetBestRegexIntValue(item.Tooltip, hitGemStatRegexes, 1)),
		proto.Stat_StatSpellCrit:  float64(GetBestRegexIntValue(item.Tooltip, critGemStatRegexes, 1)),
		proto.Stat_StatMeleeCrit:  float64(GetBestRegexIntValue(item.Tooltip, critGemStatRegexes, 1)),
		proto.Stat_StatSpellHaste: float64(GetBestRegexIntValue(item.Tooltip, hasteGemStatRegexes, 1)),
		proto.Stat_StatMeleeHaste: float64(GetBestRegexIntValue(item.Tooltip, hasteGemStatRegexes, 1)),

		proto.Stat_StatSpellPower:        float64(GetBestRegexIntValue(item.Tooltip, spellPowerGemStatRegexes, 1)),
		proto.Stat_StatHealingPower:      float64(GetBestRegexIntValue(item.Tooltip, spellPowerGemStatRegexes, 1)),
		proto.Stat_StatAttackPower:       float64(GetBestRegexIntValue(item.Tooltip, attackPowerGemStatRegexes, 1)),
		proto.Stat_StatRangedAttackPower: float64(GetBestRegexIntValue(item.Tooltip, attackPowerGemStatRegexes, 1)),
		proto.Stat_StatArmorPenetration:  float64(GetBestRegexIntValue(item.Tooltip, armorPenetrationGemStatRegexes, 1)),
		proto.Stat_StatSpellPenetration:  float64(GetBestRegexIntValue(item.Tooltip, spellPenetrationGemStatRegexes, 1)),
		proto.Stat_StatMP5:               float64(GetBestRegexIntValue(item.Tooltip, mp5GemStatRegexes, 1)),
		proto.Stat_StatExpertise:         float64(GetBestRegexIntValue(item.Tooltip, expertiseGemStatRegexes, 1)),
		proto.Stat_StatDefense:           float64(GetBestRegexIntValue(item.Tooltip, defenseGemStatRegexes, 1)),
		proto.Stat_StatDodge:             float64(GetBestRegexIntValue(item.Tooltip, dodgeGemStatRegexes, 1)),
		proto.Stat_StatParry:             float64(GetBestRegexIntValue(item.Tooltip, parryGemStatRegexes, 1)),
		proto.Stat_StatResilience:        float64(GetBestRegexIntValue(item.Tooltip, resilienceGemStatRegexes, 1)),
		proto.Stat_StatArcaneResistance:  float64(GetBestRegexIntValue(item.Tooltip, allResistGemStatRegexes, 1)),
		proto.Stat_StatFireResistance:    float64(GetBestRegexIntValue(item.Tooltip, allResistGemStatRegexes, 1)),
		proto.Stat_StatFrostResistance:   float64(GetBestRegexIntValue(item.Tooltip, allResistGemStatRegexes, 1)),
		proto.Stat_StatNatureResistance:  float64(GetBestRegexIntValue(item.Tooltip, allResistGemStatRegexes, 1)),
		proto.Stat_StatShadowResistance:  float64(GetBestRegexIntValue(item.Tooltip, allResistGemStatRegexes, 1)),
	}

	return stats
}

var itemSetNameRegex = regexp.MustCompile("<a href=\\\"\\/wotlk/item-set=-?([0-9]+)/(.*)\\\" class=\\\"q\\\">([^<]+)<")

func (item WowheadItemResponse) GetItemSetName() string {
	original := item.GetTooltipRegexString(itemSetNameRegex, 3)

	// Strip out the 10/25 man prefixes from set names
	withoutTier := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(original, "Heroes' "), "Valorous "), "Conqueror's "), "Triumphant "), "Sanctified ")
	if original != withoutTier { // if we found a tier prefix, return now.
		return withoutTier
	}

	// Now strip out the season prefix from any pvp set names
	withoutPvp := strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(original, "Savage Glad", "Glad", 1), "Hateful Glad", "Glad", 1), "Deadly Glad", "Glad", 1), "Furious Glad", "Glad", 1), "Relentless Glad", "Glad", 1), "Wrathful Glad", "Glad", 1)
	return withoutPvp
}

func getWowheadItemResponse(itemID int, tooltipsDB map[int]string) WowheadItemResponse {
	// If the db already has it, just return the db value.
	var tooltipBytes []byte

	if tooltipStr, ok := tooltipsDB[itemID]; ok {
		tooltipBytes = []byte(tooltipStr)
	} else {
		fmt.Printf("Item DB missing ID: %d\n", itemID)
		url := fmt.Sprintf("https://wowhead.com/wotlk/tooltip/item/%d?json", itemID)

		httpClient := http.Client{
			Timeout: 5 * time.Second,
		}

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err)
		}

		result, err := httpClient.Do(request)
		if err != nil {
			log.Fatal(err)
		}

		defer result.Body.Close()

		resultBody, err := ioutil.ReadAll(result.Body)
		if err != nil {
			log.Fatal(err)
		}
		tooltipBytes = resultBody
		f, err := os.OpenFile("./assets/item_data/all_item_tooltips.csv", os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("failed to open file to write: %s", err)
		}
		if strings.Contains(string(tooltipBytes), "\"error\":") {
			// fmt.Printf("Error in tooltip for %d: %s\n", i, bstr)
			log.Fatalf("failed to fetch item: %d (%s)", itemID, string(tooltipBytes))
		}
		f.WriteString(fmt.Sprintf("%d, %s, %s\n", itemID, url, tooltipBytes))
	}

	//fmt.Printf(string(tooltipStr))
	itemResponse := WowheadItemResponse{}
	err := json.Unmarshal(tooltipBytes, &itemResponse)
	if err != nil {
		fmt.Printf("Failed to decode tooltipBytes for item: %d\n", itemID)
		log.Fatal(err)
	}

	return itemResponse
}

func (item WowheadItemResponse) IsHeroic() bool {
	return strings.Contains(item.Tooltip, "<span class=\"q2\">Heroic</span>")
}

func (item WowheadItemResponse) GetRequiredProfession() proto.Profession {
	if jcGemsRegex.MatchString(item.Tooltip) {
		return proto.Profession_Jewelcrafting
	}

	return proto.Profession_ProfessionUnknown
}
