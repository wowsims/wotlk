#!/usr/bin/python

# This tool generates the talents config code, e.g. in ui/core/talents/shaman.ts.

import json
import sys

from typing import List

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from webdriver_manager.chrome import ChromeDriverManager

if len(sys.argv) < 3:
	raise Exception("Missing arguments, expected className and outputFilePath")
className = sys.argv[1]
outputFilePath = sys.argv[2]

def _between(s, start, end):
	return s[(i := s.find(start) + len(start)): i + s[i:].find(end)]


driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()))
driver.implicitly_wait(2)


def _get_spell_id_from_link(link):
	return int(link.split("/")[-2].split("=")[-1])


def get_other_spell_ranks(spell_id: int) -> List[int]:
	driver.get(f"https://wowhead.com/wotlk/spell={spell_id}#see-also-ability")

	see_also = driver.find_element(By.ID, "tab-see-also-ability")  # ordered by rank
	rows = see_also.find_elements(By.CLASS_NAME, "listview-row")
	return [_get_spell_id_from_link(row.find_element(By.CLASS_NAME, "listview-cleartext").get_attribute("href"))
		for row in rows]

def rowcol(v):
	return v["location"]["rowIdx"] + v["location"]["colIdx"]/10


to_export = []

driver.get('https://wowhead.com/wotlk/talent-calc/' + className)
trees = driver.find_elements(By.CLASS_NAME, "ctc-tree")
for tree in trees:
	_working_talents = {}

	talents = tree.find_elements(By.CLASS_NAME, "ctc-tree-talent")
	print("found %d talents\n".format(len(talents)))
	for talent in talents:
		row, col = int(talent.get_attribute("data-row")), int(talent.get_attribute("data-col"))
		max_points = int(talent.get_attribute("data-max-points"))
		link = talent.find_element(By.XPATH, "./div/a").get_attribute("href")
		name = "".join(word if i == 0 else word.title() for i, word in enumerate(link.split("/")[-1].split("-")))
		_working_talents[(row, col)] = {
			"fieldName": name,
			"location": {
				"rowIdx": row,
				"colIdx": col,
			},
			"spellIds": [_get_spell_id_from_link(link)],
			"maxPoints": max_points,
		}

	arrows = tree.find_elements(By.CLASS_NAME, "ctc-tree-talent-arrow")
	for arrow in arrows:
		prereq_row, prereq_col = int(arrow.get_attribute("data-row")), int(arrow.get_attribute("data-col"))
		length = 0
		dsstr = arrow.get_attribute("data-size")
		if dsstr:
			length = int(dsstr)

		direction = arrow.get_attribute("class").split()[-1].split("-")[-1]
		offset_row, offset_col = {"left": (0, -1), "right": (0, 1), "down": (1, 0)}[direction]

		end_row = prereq_row + offset_row * length
		end_col = prereq_col + offset_col * length

		_working_talents[(end_row, end_col)]["prereqLocation"] = {
			"rowIdx": prereq_row,
			"colIdx": prereq_col,
		}

	title = tree.find_element(By.XPATH, "./div/b").text
	background = tree.find_element(By.CLASS_NAME, "ctc-tree-talents-background").get_attribute("style")
	values = list(_working_talents.values())
	values.sort(key=rowcol)
	to_export.append({
		"name": title,
		"backgroundUrl": _between(background, '"', '"'),
		"talents": values,
	})

for subtree in to_export:
	for talent in subtree["talents"]:
		if talent["maxPoints"] > 1:
			talent["spellIds"] += get_other_spell_ranks(talent["spellIds"][0])

json_data = json.dumps(to_export, indent=2)
with open(outputFilePath, "w") as outfile:
	outfile.write(json_data)
