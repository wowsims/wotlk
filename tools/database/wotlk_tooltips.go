package database

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

type WotlkTooltipManager struct {
	TooltipManager
}

func (wtm *WotlkTooltipManager) Read() map[int32]WotlkItemResponse {
	strDB := wtm.TooltipManager.Read()
	return core.MapMap(strDB, func(id int32, tooltip string) (int32, WotlkItemResponse) {
		// Reformat the tooltip so it looks more like a wowhead tooltip.
		tooltip = strings.Replace(tooltip, fmt.Sprintf("$WowheadPower.registerItem('%d', 0, ", id), "", 1)
		tooltip = strings.TrimSuffix(tooltip, ";")
		tooltip = strings.TrimSuffix(tooltip, ")")
		tooltip = strings.ReplaceAll(tooltip, "\n", "")
		tooltip = strings.ReplaceAll(tooltip, "\t", "")
		tooltip = strings.Replace(tooltip, "name_enus: '", "\"name\": \"", 1)
		tooltip = strings.Replace(tooltip, "quality:", "\"quality\":", 1)
		tooltip = strings.Replace(tooltip, "icon: '", "\"icon\": \"", 1)
		tooltip = strings.Replace(tooltip, "tooltip_enus: '", "\"tooltip\": \"", 1)
		tooltip = strings.ReplaceAll(tooltip, "',", "\",")
		tooltip = strings.ReplaceAll(tooltip, "\\'", "'")
		// replace the '} with "}
		if strings.HasSuffix(tooltip, "'}") {
			tooltip = tooltip[:len(tooltip)-2] + "\"}"
		}

		return id, NewWotlkItemResponse(id, tooltip)
	})
}

func NewWotlkItemTooltipManager(filePath string) *WotlkTooltipManager {
	return &WotlkTooltipManager{
		TooltipManager{
			FilePath:   filePath,
			UrlPattern: "https://wotlk.evowow.com/?item=%s&power",
		},
	}
}

type WotlkItemResponse struct {
	ID      int32
	Name    string `json:"name"`
	Quality int    `json:"quality"`
	Icon    string `json:"icon"`
	Tooltip string `json:"tooltip"`
}

func NewWotlkItemResponse(id int32, tooltip string) WotlkItemResponse {
	response := WotlkItemResponse{}
	err := json.Unmarshal([]byte(tooltip), &response)
	if err != nil {
		fmt.Printf("Failed to decode tooltipBytes: %s\n", tooltip)
		panic(err)
	}
	response.ID = id
	return response
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
var wotlkdbSpellPowerRegex = regexp.MustCompile("Equip: Increases spell power by ([0-9]+)")
var wotlkdbSpellPowerRegex2 = regexp.MustCompile("Equip: Increases spell power by <!--rtg45-->([0-9]+)")

var wotlkdbHitRegex = regexp.MustCompile("Improves hit rating by <!--rtg31-->([0-9]+)")
var wotlkdbCritRegex = regexp.MustCompile("Improves critical strike rating by <!--rtg32-->([0-9]+)")
var wotlkdbHasteRegex = regexp.MustCompile("Increases your haste rating by <!--rtg36-->([0-9]+)")

var wotlkdbSpellPenetrationRegex = regexp.MustCompile("Increases your spell penetration by ([0-9]+)")
var wotlkdbMp5Regex = regexp.MustCompile("Restores ([0-9]+) mana per 5 sec")
var wotlkdbAttackPowerRegex = regexp.MustCompile(`Increases attack power by ([0-9]+)\.`)
var wotlkdbAttackPowerRegex2 = regexp.MustCompile(`Increases attack power by <!--rtg38-->([0-9]+)\.`)
var wotlkdbRangedAttackPowerRegex = regexp.MustCompile("Increases ranged attack power by ([0-9]+)")
var wotlkdbArmorPenetrationRegex = regexp.MustCompile("Increases your armor penetration rating by <!--rtg44-->([0-9]+)")
var wotlkdbArmorPenetrationRegex2 = regexp.MustCompile("Increases armor penetration rating by <!--rtg44-->([0-9]+)")
var wotlkdbExpertiseRegex = regexp.MustCompile("Increases expertise rating by <!--rtg37-->([0-9]+)")

var wotlkdbDefenseRegex = regexp.MustCompile("Equip: Increases defense rating by <!--rtg12-->([0-9]+)")
var wotlkdbDefenseRegex2 = regexp.MustCompile("Equip: Increases defense rating by ([0-9]+)")
var wotlkdbBlockRegex = regexp.MustCompile(`Equip: Increases your shield block rating by <!--rtg15-->([0-9]+)`)
var wotlkdbBlockRegex2 = regexp.MustCompile("Equip: Increases your shield block rating by ([0-9]+)")
var wotlkdbBlockValueRegex = regexp.MustCompile(`Equip: Increases the block value of your shield by ([0-9]+)\.`)
var wotlkdbBlockValueRegex2 = regexp.MustCompile("<span>([0-9]+) Block</span>")
var wotlkdbDodgeRegex = regexp.MustCompile("Increases your dodge rating by <!--rtg13-->([0-9]+)")
var wotlkdbDodgeRegex2 = regexp.MustCompile("Increases your dodge rating by ([0-9]+)")
var wotlkdbParryRegex = regexp.MustCompile("Increases your parry rating by <!--rtg14-->([0-9]+)")
var wotlkdbParryRegex2 = regexp.MustCompile("Increases your parry rating by ([0-9]+)")
var wotlkdbResilienceRegex = regexp.MustCompile("Increases your resilience rating by <!--rtg35-->([0-9]+)")
var wotlkdbArcaneResistanceRegex = regexp.MustCompile(`\+([0-9]+) Arcane Resistance`)
var wotlkdbFireResistanceRegex = regexp.MustCompile(`\+([0-9]+) Fire Resistance`)
var wotlkdbFrostResistanceRegex = regexp.MustCompile(`\+([0-9]+) Frost Resistance`)
var wotlkdbNatureResistanceRegex = regexp.MustCompile(`\+([0-9]+) Nature Resistance`)
var wotlkdbShadowResistanceRegex = regexp.MustCompile(`\+([0-9]+) Shadow Resistance`)

func (item WotlkItemResponse) GetStats() Stats {
	sp := float64(item.GetIntValue(wotlkdbSpellPowerRegex)) + float64(item.GetIntValue(wotlkdbSpellPowerRegex2))
	return Stats{
		proto.Stat_StatArmor:             float64(item.GetIntValue(wotlkdbArmorRegex)),
		proto.Stat_StatStrength:          float64(item.GetIntValue(wotlkdbStrengthRegex)),
		proto.Stat_StatAgility:           float64(item.GetIntValue(wotlkdbAgilityRegex)),
		proto.Stat_StatStamina:           float64(item.GetIntValue(wotlkdbStaminaRegex)),
		proto.Stat_StatIntellect:         float64(item.GetIntValue(wotlkdbIntellectRegex)),
		proto.Stat_StatSpirit:            float64(item.GetIntValue(wotlkdbSpiritRegex)),
		proto.Stat_StatSpellPower:        sp,
		proto.Stat_StatSpellHit:          float64(item.GetIntValue(wotlkdbHitRegex)),
		proto.Stat_StatMeleeHit:          float64(item.GetIntValue(wotlkdbHitRegex)),
		proto.Stat_StatSpellCrit:         float64(item.GetIntValue(wotlkdbCritRegex)),
		proto.Stat_StatMeleeCrit:         float64(item.GetIntValue(wotlkdbCritRegex)),
		proto.Stat_StatSpellHaste:        float64(item.GetIntValue(wotlkdbHasteRegex)),
		proto.Stat_StatMeleeHaste:        float64(item.GetIntValue(wotlkdbHasteRegex)),
		proto.Stat_StatSpellPenetration:  float64(item.GetIntValue(wotlkdbSpellPenetrationRegex)),
		proto.Stat_StatMP5:               float64(item.GetIntValue(wotlkdbMp5Regex)),
		proto.Stat_StatAttackPower:       float64(item.GetIntValue(wotlkdbAttackPowerRegex) + item.GetIntValue(wotlkdbAttackPowerRegex2)),
		proto.Stat_StatRangedAttackPower: float64(item.GetIntValue(wotlkdbAttackPowerRegex) + item.GetIntValue(wotlkdbAttackPowerRegex2) + item.GetIntValue(wotlkdbRangedAttackPowerRegex)),
		proto.Stat_StatArmorPenetration:  float64(item.GetIntValue(wotlkdbArmorPenetrationRegex) + item.GetIntValue(wotlkdbArmorPenetrationRegex2)),
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

func (item WotlkItemResponse) IsPattern() bool {
	for _, pattern := range patternRegexes {
		if pattern.MatchString(item.Tooltip) {
			return true
		}
	}
	return false
}

func (item WotlkItemResponse) IsRandomEnchant() bool {
	return randomEnchantRegex.MatchString(item.Tooltip)
}

func (item WotlkItemResponse) IsEquippable() bool {
	return item.GetItemType() != proto.ItemType_ItemTypeUnknown &&
		!item.IsPattern() &&
		!item.IsRandomEnchant()
}

var wotlkItemLevelRegex = regexp.MustCompile("Item Level ([0-9]+)<")

func (item WotlkItemResponse) GetItemLevel() int {
	return item.GetIntValue(wotlkItemLevelRegex)
}

// WOTLK DB has no phase info
func (item WotlkItemResponse) GetPhase() int {

	ilvl := item.GetItemLevel()
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

func (item WotlkItemResponse) GetUnique() bool {
	return uniqueRegex.MatchString(item.Tooltip) && !jcGemsRegex.MatchString(item.Tooltip)
}

func (item WotlkItemResponse) GetItemType() proto.ItemType {
	for itemType, pattern := range itemTypePatterns {
		if pattern.MatchString(item.Tooltip) {
			return itemType
		}
	}
	return proto.ItemType_ItemTypeUnknown
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
	proto.RangedWeaponType_RangedWeaponTypeIdol:     regexp.MustCompile("<th><!--asc8-->Idol</th>"),
	proto.RangedWeaponType_RangedWeaponTypeLibram:   regexp.MustCompile("<th><!--asc7-->Libram</th>"),
	proto.RangedWeaponType_RangedWeaponTypeThrown:   regexp.MustCompile("<th>Thrown</th>"),
	proto.RangedWeaponType_RangedWeaponTypeTotem:    regexp.MustCompile("<th><!--asc9-->Totem</th>"),
	proto.RangedWeaponType_RangedWeaponTypeWand:     regexp.MustCompile("<th>Wand</th>"),
	proto.RangedWeaponType_RangedWeaponTypeSigil:    regexp.MustCompile("<th><!--asc10-->Sigil</th>"),
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
	} else if matches := weaponDamageRegex2.FindStringSubmatch(item.Tooltip); len(matches) > 0 {
		val, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse weapon damage: %s", err)
		}
		return val, val
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
		proto.Stat_StatSpellHaste:        float64(GetBestRegexIntValue(bonusStr, hasteSocketBonusRegexes, 1)),
		proto.Stat_StatSpellPower:        float64(GetBestRegexIntValue(bonusStr, spellPowerSocketBonusRegexes, 1)),
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

var wotlkItemSetNameRegex = regexp.MustCompile("<a href=\\\"\\?itemset=([0-9]+)\\\" class=\\\"q\\\">([^<]+)<")

func (item WotlkItemResponse) GetItemSetName() string {
	return item.GetTooltipRegexString(wotlkItemSetNameRegex, 2)
}

func (item WotlkItemResponse) IsHeroic() bool {
	return strings.Contains(item.Tooltip, "<span class=\"q2\">Heroic</span>")
}

func (item WotlkItemResponse) GetRequiredProfession() proto.Profession {
	if jcGemsRegex.MatchString(item.Tooltip) {
		return proto.Profession_Jewelcrafting
	}

	return proto.Profession_ProfessionUnknown
}
