#!/usr/bin/python

# Usage example:
# python3 ./tools/scrape_enchant_descriptions.py ./sim/core/items/all_enchants.go ./assets/enchants/descriptions.json

import json
import re
import sys

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from selenium.webdriver.common.desired_capabilities import DesiredCapabilities
from selenium.webdriver.common.keys import Keys
from webdriver_manager.chrome import ChromeDriverManager

if len(sys.argv) < 3:
	raise Exception("Missing arguments, expected input_file_path and output_file_path")
input_file_path = sys.argv[1]
output_file_path = sys.argv[2]

input_file = open("sim/core/items/all_enchants.go", 'r')
input_lines = input_file.readlines()

enchants = []
for line in input_lines:
	id_match = re.search(r"ID:\s*(\d+)", line)
	if id_match is None:
		continue
	enchant_id = int(id_match.group(1))

	effect_id_match = re.search(r"EffectID:\s*(\d+)", line)
	effect_id = int(effect_id_match.group(1))

	is_spell_id = "IsSpellID" in line

	enchants.append({
		"id": enchant_id,
		"effect_id": effect_id,
		"is_spell_id": is_spell_id,
	})

def get_spell_url(item_id):
	driver.get("https://wowhead.com/wotlk/item=" + str(item_id))
	tooltips = driver.find_elements(By.CLASS_NAME, "wowhead-tooltip")
	tooltip = tooltips[0]
	anchors = tooltip.find_elements(By.TAG_NAME, "a")
	for anchor in anchors:
		href = anchor.get_attribute("href")
		if "/wotlk/spell=" in href:
			print("Item {} has spell url {}\n".format(item_id, href))
			return href
	raise Exception("No results for id " + str(item_id))

def get_spell_effect_description(spell_url):
	driver.get(spell_url)
	details_table = driver.find_elements(By.ID, "spelldetails")[0]
	effect_elem = details_table.find_elements(By.CLASS_NAME, "q2")[0]
	print("Spell {} has description {}".format(spell_url, effect_elem.text))
	return effect_elem.text

def get_enchant_description(enchant):
	if enchant["is_spell_id"]:
		return get_spell_effect_description("https://wowhead.com/wotlk/spell={}".format(enchant["id"]))
	else:
		spell_url = get_spell_url(enchant["id"])
		return get_spell_effect_description(spell_url)

caps = DesiredCapabilities().CHROME
caps["pageLoadStrategy"] = "eager"   # Do not wait for full page load
driver = webdriver.Chrome(desired_capabilities=caps, service=Service(ChromeDriverManager().install()))
driver.implicitly_wait(2)

for enchant in enchants:
	enchant["description"] = get_enchant_description(enchant)

driver.quit()

with open(output_file_path, "w") as outfile:
	outfile.write("{\n")
	for i, enchant in enumerate(enchants):
		outfile.write("\t\"{}\": \"{}\"{}\n".format(enchant["effect_id"], enchant["description"], "" if i == len(enchants) - 1 else ","))
	outfile.write("}")
