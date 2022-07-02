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

func (item WotlkItemResponse) GetStats() Stats {
	return Stats{
		proto.Stat_StatArmor:             float64(item.GetIntValue(armorRegex)),
		proto.Stat_StatStrength:          float64(item.GetIntValue(strengthRegex)),
		proto.Stat_StatAgility:           float64(item.GetIntValue(agilityRegex)),
		proto.Stat_StatStamina:           float64(item.GetIntValue(staminaRegex)),
		proto.Stat_StatIntellect:         float64(item.GetIntValue(intellectRegex)),
		proto.Stat_StatSpirit:            float64(item.GetIntValue(spiritRegex)),
		proto.Stat_StatSpellPower:        float64(item.GetIntValue(spellPowerRegex)),
		proto.Stat_StatHealingPower:      float64(item.GetIntValue(spellPowerRegex)),
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
		proto.Stat_StatAttackPower:       float64(item.GetIntValue(attackPowerRegex)),
		proto.Stat_StatRangedAttackPower: float64(item.GetIntValue(attackPowerRegex) + item.GetIntValue(rangedAttackPowerRegex)),
		proto.Stat_StatFeralAttackPower:  float64(item.GetIntValue(feralAttackPowerRegex)),
		proto.Stat_StatArmorPenetration:  float64(item.GetIntValue(armorPenetrationRegex)),
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
	proto.WeaponType_WeaponTypeShield:  regexp.MustCompile("<th>Shield</th>"),
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

var wotlkItemSetNameRegex = regexp.MustCompile("<a href=\\\"\\/itemset=([0-9]+)\\\" class=\\\"q\\\">([^<]+)<")

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
