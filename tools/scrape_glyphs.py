#!/usr/bin/python

# This tool generates the glyphs proto code and UI config code, e.g.
# 'ShamanMajorGlyphs' and 'ShamanMinorGlyphs' in proto/shaman.proto
# and the config in ui/core/talents/shaman.ts.

import json
import re
import sys

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from webdriver_manager.chrome import ChromeDriverManager

if len(sys.argv) < 3:
	raise Exception("Missing arguments, expected class_name and output_file_path")
class_name = sys.argv[1]
output_file_path = sys.argv[2]

# Convert "death-knight" to DeathKnight
pretty_class_name = "".join(word.title() for i, word in enumerate(class_name.split("-")))
lower_class_name = "".join(word if i == 0 else word.title() for i, word in enumerate(class_name.split("-")))

def get_glyphs_data(glyph_button):
	glyph_button.click()
	glyphs_menu = driver.find_element(By.CLASS_NAME, "ctc-glyphs-picker-listview")
	menu_rows = glyphs_menu.find_elements(By.CLASS_NAME, "listview-row")
	glyphs_data = []

	for menu_row in menu_rows:
		cells = menu_row.find_elements(By.TAG_NAME, "td")
		label_elem = cells[1].find_element(By.TAG_NAME, "a")
		glyph_name = label_elem.text
		print("Glyph name: " + glyph_name)
		if glyph_name == "None":
			continue

		link = label_elem.get_attribute("href")
		id_match = re.search(r"item=(\d+)", link)
		glyph_id = int(id_match.group(1))

		ins_elem = cells[0].find_element(By.TAG_NAME, "ins")
		bg_style = ins_elem.get_attribute("style")
		icon_url = re.search(r"url\(\"(.*)\"\)", bg_style).group(1)
		icon_url = icon_url.replace("icons/small", "icons/large")

		glyphs_data.append({
			"name": glyph_name,
			"id": glyph_id,
			"description": cells[2].text,
			"icon_url": icon_url,
		})

	return glyphs_data

driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()))
driver.implicitly_wait(2)

driver.get('https://wowhead.com/wotlk/talent-calc/' + class_name)
glyph_slots = driver.find_elements(By.CLASS_NAME, "ctc-glyphs-group-slot")
major_glyph_slot = next(gs for gs in glyph_slots if int(gs.get_attribute("data-glyph-slot")) == 0)
minor_glyph_slot = next(gs for gs in glyph_slots if int(gs.get_attribute("data-glyph-slot")) == 3)

major_glyphs_data = get_glyphs_data(major_glyph_slot)
webdriver.ActionChains(driver).send_keys(Keys.ESCAPE).perform()
minor_glyphs_data = get_glyphs_data(minor_glyph_slot)
driver.quit()

def write_glyphs_proto(outfile, glyphs_data, glyph_type):
	outfile.write("enum {}{}Glyph {{\n".format(pretty_class_name, glyph_type))
	outfile.write("\t{}{}GlyphNone = 0;\n".format(pretty_class_name, glyph_type))

	glyph_idx = 1
	for glyph_data in glyphs_data:
		pretty_glyph_name = re.sub(r"\W+", "", glyph_data["name"].title())
		outfile.write("\t{} = {};\n".format(pretty_glyph_name, glyph_data["id"]))

	outfile.write("}\n")

def write_glyphs_config(outfile, glyphs_data, glyph_type):
	outfile.write("\t{}Glyphs: {{\n".format(glyph_type.lower()))

	for glyph_data in glyphs_data:
		pretty_glyph_name = re.sub(r"\W+", "", glyph_data["name"].title())
		outfile.write("\t\t[{}{}Glyph.{}]: {{\n".format(pretty_class_name, glyph_type, pretty_glyph_name, glyph_data["icon_url"]))
		outfile.write("\t\t\tname: '{}',\n".format(glyph_data["name"].replace("'", "\\'")))
		outfile.write("\t\t\tdescription: '{}',\n".format(glyph_data["description"].replace("'", "\\'")))
		outfile.write("\t\t\ticonUrl: '{}',\n".format(glyph_data["icon_url"]))
		outfile.write("\t\t},\n")

	outfile.write("\t},\n")


with open(output_file_path, "w") as outfile:
	write_glyphs_proto(outfile, major_glyphs_data, "Major")
	write_glyphs_proto(outfile, minor_glyphs_data, "Minor")

	outfile.write("\n")
	outfile.write("export const {}GlyphsConfig: GlyphsConfig = {{\n".format(lower_class_name))
	write_glyphs_config(outfile, major_glyphs_data, "Major")
	write_glyphs_config(outfile, minor_glyphs_data, "Minor")
	outfile.write("};")
