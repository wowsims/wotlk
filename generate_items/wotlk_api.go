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

type WotlkItemResponse struct {
	Name    string `json:"name"`
	Quality int    `json:"quality"`
	Icon    string `json:"icon"`
	Tooltip string `json:"tooltip"`
}

func (item WotlkItemResponse) GetName() string {
	return item.Name
}
func (item WotlkItemResponse) GetQuality() int {
	return item.Quality
}
func (item WotlkItemResponse) GetIcon() string {
	return item.Icon
}

func (item WotlkItemResponse) TooltipWithoutSetBonus() string {
	setIdx := strings.Index(item.Tooltip, "Set : ")
	if setIdx == -1 {
		return item.Tooltip
	} else {
		return item.Tooltip[:setIdx]
	}
}

func (item WotlkItemResponse) GetTooltipRegexString(pattern *regexp.Regexp, matchIdx int) string {
	return GetRegexStringValue(item.TooltipWithoutSetBonus(), pattern, matchIdx)
}

func (item WotlkItemResponse) GetTooltipRegexValue(pattern *regexp.Regexp, matchIdx int) int {
	return GetRegexIntValue(item.TooltipWithoutSetBonus(), pattern, matchIdx)
}

func (item WotlkItemResponse) GetIntValue(pattern *regexp.Regexp) int {
	return item.GetTooltipRegexValue(pattern, 1)
}

var wotlkdbArmorRegex = regexp.MustCompile("<!--amr-->([0-9]+) Armor")
var wotlkdbAgilityRegex = regexp.MustCompile(`<!--stat3-->\+([0-9]+) Agility`)
var wotlkdbStrengthRegex = regexp.MustCompile(`<!--stat4-->\+([0-9]+) Strength`)
var wotlkdbIntellectRegex = regexp.MustCompile(`<!--stat5-->\+([0-9]+) Intellect`)
var wotlkdbSpiritRegex = regexp.MustCompile(`<!--stat6-->\+([0-9]+) Spirit`)
var wotlkdbStaminaRegex = regexp.MustCompile(`<!--stat7-->\+([0-9]+) Stamina`)
var wotlkdbSpellPowerRegex = regexp.MustCompile("Increases spell power by ([0-9]+)")

// Not sure these exist anymore?
var wotlkdbArcaneSpellPowerRegex = regexp.MustCompile("Increases Arcane power by ([0-9]+)")
var wotlkdbFireSpellPowerRegex = regexp.MustCompile("Increases Fire power by ([0-9]+)")
var wotlkdbFrostSpellPowerRegex = regexp.MustCompile("Increases Frost power by ([0-9]+)")
var wotlkdbHolySpellPowerRegex = regexp.MustCompile("Increases Holy power by ([0-9]+)")
var wotlkdbNatureSpellPowerRegex = regexp.MustCompile("Increases Nature power by ([0-9]+)")
var wotlkdbShadowSpellPowerRegex = regexp.MustCompile("Increases Shadow power by ([0-9]+)")

var wotlkdbHitRegex = regexp.MustCompile("Improves hit rating by <!--rtg31-->([0-9]+)")
var wotlkdbCritRegex = regexp.MustCompile("Improves critical strike rating by <!--rtg32-->([0-9]+)")
var wotlkdbHasteRegex = regexp.MustCompile("Increases your haste rating by <!--rtg36-->([0-9]+)")

var wotlkdbSpellPenetrationRegex = regexp.MustCompile("Increases your spell penetration by ([0-9]+)")
var wotlkdbMp5Regex = regexp.MustCompile("Restores ([0-9]+) mana per 5 sec")
var wotlkdbAttackPowerRegex = regexp.MustCompile(`Increases attack power by ([0-9]+)\.`)
var wotlkdbRangedAttackPowerRegex = regexp.MustCompile("Increases ranged attack power by ([0-9]+)")
var wotlkdbArmorPenetrationRegex = regexp.MustCompile("Increases armor penetration rating by <!--rtg44-->([0-9]+)")
var wotlkdbExpertiseRegex = regexp.MustCompile("Increases expertise rating by <!--rtg37-->([0-9]+)")

var wotlkdbDefenseRegex = regexp.MustCompile("Equip: Increases defense rating by <!--rtg12-->([0-9]+)")
var wotlkdbDefenseRegex2 = regexp.MustCompile("Equip: Increases defense rating by ([0-9]+)")
var wotlkdbBlockRegex = regexp.MustCompile(`Equip: Increases your shield block rating by <!--rtg15-->([0-9]+)\.`)
var wotlkdbBlockRegex2 = regexp.MustCompile("Equip: Increases your shield block rating by ([0-9]+)")
var wotlkdbBlockValueRegex = regexp.MustCompile(`Equip: Increases the block value of your shield by ([0-9]+)\.`)
var wotlkdbBlockValueRegex2 = regexp.MustCompile("<span>([0-9]+) Block</span>")
var wotlkdbDodgeRegex = regexp.MustCompile("Increases your dodge rating by <!--rtg13-->([0-9]+)")
var wotlkdbDodgeRegex2 = regexp.MustCompile("Increases your dodge rating by ([0-9]+)")
var wotlkdbParryRegex = regexp.MustCompile("Increases your parry rating by <!--rtg14-->([0-9]+)")
var wotlkdbParryRegex2 = regexp.MustCompile("Increases your parry rating by ([0-9]+)")
var wotlkdbResilienceRegex = regexp.MustCompile("Improves your resilience rating by <!--rtg35-->([0-9]+)")
var wotlkdbArcaneResistanceRegex = regexp.MustCompile(`\+([0-9]+) Arcane Resistance`)
var wotlkdbFireResistanceRegex = regexp.MustCompile(`\+([0-9]+) Fire Resistance`)
var wotlkdbFrostResistanceRegex = regexp.MustCompile(`\+([0-9]+) Frost Resistance`)
var wotlkdbNatureResistanceRegex = regexp.MustCompile(`\+([0-9]+) Nature Resistance`)
var wotlkdbShadowResistanceRegex = regexp.MustCompile(`\+([0-9]+) Shadow Resistance`)

func (item WotlkItemResponse) GetStats() Stats {
	return Stats{
		proto.Stat_StatArmor:             float64(item.GetIntValue(wotlkdbArmorRegex)),
		proto.Stat_StatStrength:          float64(item.GetIntValue(wotlkdbStrengthRegex)),
		proto.Stat_StatAgility:           float64(item.GetIntValue(wotlkdbAgilityRegex)),
		proto.Stat_StatStamina:           float64(item.GetIntValue(wotlkdbStaminaRegex)),
		proto.Stat_StatIntellect:         float64(item.GetIntValue(wotlkdbIntellectRegex)),
		proto.Stat_StatSpirit:            float64(item.GetIntValue(wotlkdbSpiritRegex)),
		proto.Stat_StatSpellPower:        float64(item.GetIntValue(wotlkdbSpellPowerRegex)),
		proto.Stat_StatHealingPower:      float64(item.GetIntValue(wotlkdbSpellPowerRegex)),
		proto.Stat_StatArcaneSpellPower:  float64(item.GetIntValue(wotlkdbArcaneSpellPowerRegex)),
		proto.Stat_StatFireSpellPower:    float64(item.GetIntValue(wotlkdbFireSpellPowerRegex)),
		proto.Stat_StatFrostSpellPower:   float64(item.GetIntValue(wotlkdbFrostSpellPowerRegex)),
		proto.Stat_StatHolySpellPower:    float64(item.GetIntValue(wotlkdbHolySpellPowerRegex)),
		proto.Stat_StatNatureSpellPower:  float64(item.GetIntValue(wotlkdbNatureSpellPowerRegex)),
		proto.Stat_StatShadowSpellPower:  float64(item.GetIntValue(wotlkdbShadowSpellPowerRegex)),
		proto.Stat_StatSpellHit:          float64(item.GetIntValue(wotlkdbHitRegex)),
		proto.Stat_StatMeleeHit:          float64(item.GetIntValue(wotlkdbHitRegex)),
		proto.Stat_StatSpellCrit:         float64(item.GetIntValue(wotlkdbCritRegex)),
		proto.Stat_StatMeleeCrit:         float64(item.GetIntValue(wotlkdbCritRegex)),
		proto.Stat_StatSpellHaste:        float64(item.GetIntValue(wotlkdbHasteRegex)),
		proto.Stat_StatMeleeHaste:        float64(item.GetIntValue(wotlkdbHasteRegex)),
		proto.Stat_StatSpellPenetration:  float64(item.GetIntValue(wotlkdbSpellPenetrationRegex)),
		proto.Stat_StatMP5:               float64(item.GetIntValue(wotlkdbMp5Regex)),
		proto.Stat_StatAttackPower:       float64(item.GetIntValue(wotlkdbAttackPowerRegex)),
		proto.Stat_StatRangedAttackPower: float64(item.GetIntValue(wotlkdbAttackPowerRegex) + item.GetIntValue(wotlkdbRangedAttackPowerRegex)),
		proto.Stat_StatArmorPenetration:  float64(item.GetIntValue(wotlkdbArmorPenetrationRegex)),
		proto.Stat_StatExpertise:         float64(item.GetIntValue(wotlkdbExpertiseRegex)),
		proto.Stat_StatDefense:           float64(item.GetIntValue(wotlkdbDefenseRegex) + item.GetIntValue(wotlkdbDefenseRegex2)),
		proto.Stat_StatBlock:             float64(item.GetIntValue(wotlkdbBlockRegex) + item.GetIntValue(wotlkdbBlockRegex2)),
		proto.Stat_StatBlockValue:        float64(item.GetIntValue(wotlkdbBlockValueRegex) + item.GetIntValue(wotlkdbBlockValueRegex2)),
		proto.Stat_StatDodge:             float64(item.GetIntValue(wotlkdbDodgeRegex) + item.GetIntValue(wotlkdbDodgeRegex2)),
		proto.Stat_StatParry:             float64(item.GetIntValue(wotlkdbParryRegex) + item.GetIntValue(wotlkdbParryRegex2)),
		proto.Stat_StatResilience:        float64(item.GetIntValue(wotlkdbResilienceRegex)),
		proto.Stat_StatArcaneResistance:  float64(item.GetIntValue(wotlkdbArcaneResistanceRegex)),
		proto.Stat_StatFireResistance:    float64(item.GetIntValue(wotlkdbFireResistanceRegex)),
		proto.Stat_StatFrostResistance:   float64(item.GetIntValue(wotlkdbFrostResistanceRegex)),
		proto.Stat_StatNatureResistance:  float64(item.GetIntValue(wotlkdbNatureResistanceRegex)),
		proto.Stat_StatShadowResistance:  float64(item.GetIntValue(wotlkdbShadowResistanceRegex)),
	}
}

func (item WotlkItemResponse) GetClassAllowlist() []proto.Class {
	var allowlist []proto.Class

	for _, entry := range classPatterns {
		if entry.pattern.MatchString(item.Tooltip) {
			allowlist = append(allowlist, entry.class)
		}
	}

	return allowlist
}

func (item WotlkItemResponse) IsEquippable() bool {
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

func (item WotlkItemResponse) IsPattern() bool {
	for _, pattern := range nonEquippableRegexes {
		if pattern.MatchString(item.Tooltip) {
			return true
		}
	}
	return false
}

var wotlkItemLevelRegex = regexp.MustCompile("Item Level ([0-9]+)<")

func (item WotlkItemResponse) GetItemLevel() int {
	return item.GetIntValue(wotlkItemLevelRegex)
}

// WOTLK DB has no phase info
func (item WotlkItemResponse) GetPhase() int {
	return 1
}

func (item WotlkItemResponse) GetUnique() bool {
	return uniqueRegex.MatchString(item.Tooltip)
}

func (item WotlkItemResponse) GetItemType() proto.ItemType {
	for itemType, pattern := range itemTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return itemType
		}
	}
	panic("Could not find item type from tooltip: " + item.Tooltip)
}

var wotlkArmorTypePatterns = map[proto.ArmorType]*regexp.Regexp{
	proto.ArmorType_ArmorTypeCloth:   regexp.MustCompile("<th><!--asc1-->Cloth</th>"),
	proto.ArmorType_ArmorTypeLeather: regexp.MustCompile("<th><!--asc2-->Leather</th>"),
	proto.ArmorType_ArmorTypeMail:    regexp.MustCompile("<th><!--asc3-->Mail</th>"),
	proto.ArmorType_ArmorTypePlate:   regexp.MustCompile("<th><!--asc4-->Plate</th>"),
}

func (item WotlkItemResponse) GetArmorType() proto.ArmorType {
	for armorType, pattern := range wotlkArmorTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return armorType
		}
	}
	return proto.ArmorType_ArmorTypeUnknown
}

var wotlkWeaponTypePatterns = map[proto.WeaponType]*regexp.Regexp{
	proto.WeaponType_WeaponTypeAxe:     regexp.MustCompile("<th>Axe</th>"),
	proto.WeaponType_WeaponTypeDagger:  regexp.MustCompile("<th>Dagger</th>"),
	proto.WeaponType_WeaponTypeFist:    regexp.MustCompile("<th>Fist Weapon</th>"),
	proto.WeaponType_WeaponTypeMace:    regexp.MustCompile("<th>Mace</th>"),
	proto.WeaponType_WeaponTypeOffHand: regexp.MustCompile("<td>Held In Off-Hand</td>"),
	proto.WeaponType_WeaponTypePolearm: regexp.MustCompile("<th>Polearm</th>"),
	proto.WeaponType_WeaponTypeShield:  regexp.MustCompile("<th><!--asc6-->Shield</th>"),
	proto.WeaponType_WeaponTypeStaff:   regexp.MustCompile("<th>Staff</th>"),
	proto.WeaponType_WeaponTypeSword:   regexp.MustCompile("<th>Sword</th>"),
}

func (item WotlkItemResponse) GetWeaponType() proto.WeaponType {
	for weaponType, pattern := range wotlkWeaponTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return weaponType
		}
	}
	return proto.WeaponType_WeaponTypeUnknown
}

func (item WotlkItemResponse) GetHandType() proto.HandType {
	for handType, pattern := range handTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return handType
		}
	}
	return proto.HandType_HandTypeUnknown
}

var wotlkRangedWeaponTypePatterns = map[proto.RangedWeaponType]*regexp.Regexp{
	proto.RangedWeaponType_RangedWeaponTypeBow:      regexp.MustCompile("<th>Bow</th>"),
	proto.RangedWeaponType_RangedWeaponTypeCrossbow: regexp.MustCompile("<th>Crossbow</th>"),
	proto.RangedWeaponType_RangedWeaponTypeGun:      regexp.MustCompile("<th>Gun</th>"),
	proto.RangedWeaponType_RangedWeaponTypeIdol:     regexp.MustCompile("<th>Idol</th>"),
	proto.RangedWeaponType_RangedWeaponTypeLibram:   regexp.MustCompile("<th>Libram</th>"),
	proto.RangedWeaponType_RangedWeaponTypeThrown:   regexp.MustCompile("<th>Thrown</th>"),
	proto.RangedWeaponType_RangedWeaponTypeTotem:    regexp.MustCompile("<th>Totem</th>"),
	proto.RangedWeaponType_RangedWeaponTypeWand:     regexp.MustCompile("<th>Wand</th>"),
}

func (item WotlkItemResponse) GetRangedWeaponType() proto.RangedWeaponType {
	for rangedWeaponType, pattern := range wotlkRangedWeaponTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return rangedWeaponType
		}
	}
	return proto.RangedWeaponType_RangedWeaponTypeUnknown
}

// Returns min/max of weapon damage
func (item WotlkItemResponse) GetWeaponDamage() (float64, float64) {
	if matches := weaponDamageRegex.FindStringSubmatch(item.Tooltip); len(matches) > 0 {
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
	}
	return 0, 0
}

func (item WotlkItemResponse) GetWeaponSpeed() float64 {
	if matches := weaponSpeedRegex.FindStringSubmatch(item.Tooltip); len(matches) > 0 {
		speed, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		return speed
	}
	return 0
}

func (item WotlkItemResponse) GetGemSockets() []proto.GemColor {
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

func (item WotlkItemResponse) GetSocketBonus() Stats {
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
		proto.Stat_StatSpellPower:        float64(GetBestRegexIntValue(bonusStr, spellPowerSocketBonusRegexes, 1)),
		proto.Stat_StatHealingPower:      float64(GetBestRegexIntValue(bonusStr, spellPowerSocketBonusRegexes, 1)),
		proto.Stat_StatSpellHit:          float64(GetBestRegexIntValue(bonusStr, spellHitSocketBonusRegexes, 1)),
		proto.Stat_StatMeleeHit:          float64(GetBestRegexIntValue(bonusStr, spellHitSocketBonusRegexes, 1)),
		proto.Stat_StatSpellCrit:         float64(GetBestRegexIntValue(bonusStr, spellCritSocketBonusRegexes, 1)),
		proto.Stat_StatMeleeCrit:         float64(GetBestRegexIntValue(bonusStr, spellCritSocketBonusRegexes, 1)),
		proto.Stat_StatMP5:               float64(GetBestRegexIntValue(bonusStr, mp5SocketBonusRegexes, 1)),
		proto.Stat_StatAttackPower:       float64(GetBestRegexIntValue(bonusStr, attackPowerSocketBonusRegexes, 1)),
		proto.Stat_StatRangedAttackPower: float64(GetBestRegexIntValue(bonusStr, attackPowerSocketBonusRegexes, 1)),
		proto.Stat_StatDefense:           float64(GetBestRegexIntValue(bonusStr, defenseSocketBonusRegexes, 1)),
		proto.Stat_StatBlock:             float64(GetBestRegexIntValue(bonusStr, blockSocketBonusRegexes, 1)),
		proto.Stat_StatDodge:             float64(GetBestRegexIntValue(bonusStr, dodgeSocketBonusRegexes, 1)),
		proto.Stat_StatParry:             float64(GetBestRegexIntValue(bonusStr, parrySocketBonusRegexes, 1)),
		proto.Stat_StatResilience:        float64(GetBestRegexIntValue(bonusStr, resilienceSocketBonusRegexes, 1)),
	}

	return stats
}

func (item WotlkItemResponse) GetSocketColor() proto.GemColor {
	for socketColor, pattern := range gemSocketColorPatterns {
		if pattern.MatchString(item.Tooltip) {
			return socketColor
		}
	}
	// fmt.Printf("Could not find socket color for gem %s\n", item.Name)
	return proto.GemColor_GemColorUnknown
}

func (item WotlkItemResponse) GetGemStats() Stats {
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
		proto.Stat_StatSpellPenetration:  float64(GetBestRegexIntValue(item.Tooltip, spellPenetrationGemStatRegexes, 1)),
		proto.Stat_StatMP5:               float64(GetBestRegexIntValue(item.Tooltip, mp5GemStatRegexes, 1)),
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

var wotlkItemSetNameRegex = regexp.MustCompile("<a href=\\\"\\?itemset=([0-9]+)\\\" class=\\\"q\\\">([^<]+)<")

func (item WotlkItemResponse) GetItemSetName() string {
	return item.GetTooltipRegexString(wotlkItemSetNameRegex, 2)
}

func getWotlkItemResponse(itemID int, tooltipsDB map[int]string) WotlkItemResponse {
	// If the db already has it, just return the db value.
	var tooltipBytes []byte

	if tooltipStr, ok := tooltipsDB[itemID]; ok {
		tooltipBytes = []byte(tooltipStr)
	} else {
		fmt.Printf("Item DB missing ID: %d\n", itemID)
		url := fmt.Sprintf("https://wotlkdb.com/?item=%d&power", itemID)

		httpClient := http.Client{
			Timeout: 5 * time.Second,
		}

		result, err := httpClient.Get(url)
		if err != nil {
			fmt.Printf("Error fetching %d: %s\n", itemID, err)
			return WotlkItemResponse{}
		}
		defer result.Body.Close()

		body, _ := ioutil.ReadAll(result.Body)
		bstr := string(body)
		bstr = strings.Replace(bstr, fmt.Sprintf("$WowheadPower.registerItem('%d', 0, ", itemID), "", 1)
		bstr = strings.TrimSuffix(bstr, ";")
		bstr = strings.TrimSuffix(bstr, ")")
		bstr = strings.ReplaceAll(bstr, "\n", "")
		bstr = strings.ReplaceAll(bstr, "\t", "")
		bstr = strings.Replace(bstr, "name_enus: '", "\"name\": \"", 1)
		bstr = strings.Replace(bstr, "quality:", "\"quality\":", 1)
		bstr = strings.Replace(bstr, "icon: '", "\"icon\": \"", 1)
		bstr = strings.Replace(bstr, "tooltip_enus: '", "\"tooltip\": \"", 1)
		bstr = strings.ReplaceAll(bstr, "',", "\",")
		bstr = strings.ReplaceAll(bstr, "\\'", "'")
		// replace the '} with "}
		if strings.HasSuffix(bstr, "'}") {
			bstr = bstr[:len(bstr)-2] + "\"}"
		}

		fmt.Printf("Found Item %d: %s\n", itemID, bstr)

		fmt.Printf("Writing to all_item_tooltips.csv now...\n")
		tooltipBytes = []byte(bstr)
		alltooltips, err := os.OpenFile("./assets/item_data/all_item_tooltips.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		line := fmt.Sprintf("%d, %s, %s\n", itemID, url, bstr)
		n, err := alltooltips.WriteString(line)
		if n != len(line) {
			log.Fatalf("Unexpected number of bytes written: %d ... expected: %d", n, len(line))
		}
		if err != nil {
			log.Fatalf("Failed to append to item tooltips: %s", err)
		}
		err = alltooltips.Close()
		if err != nil {
			log.Fatalf("Failed to close item tooltips: %s", err)
		}
	}

	//fmt.Printf(string(tooltipStr))
	itemResponse := WotlkItemResponse{}
	err := json.Unmarshal(tooltipBytes, &itemResponse)
	if err != nil {
		fmt.Printf("Failed to decode tooltipBytes for item: %d\n", itemID)
		log.Fatal(err)
	}

	return itemResponse
}

func (item WotlkItemResponse) IsHeroic() bool {
	return strings.Contains(item.Tooltip, "<span class=\"q2\">Heroic</span>")
}
