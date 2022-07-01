#!/usr/bin/python

import json
import sys

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from webdriver_manager.chrome import ChromeDriverManager

if len(sys.argv) < 3:
	raise Exception("Missing arguments, expected className and outputFilePath")
className = sys.argv[1]
outputFilePath = sys.argv[2]

driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()))
driver.implicitly_wait(2)

driver.get('https://wowhead.com/wotlk/talent-calc/' + className)
trees = driver.find_elements(By.CLASS_NAME, "ctc-tree")

with open(outputFilePath, "w") as outfile:
	# Convert "death-knight" to DeathKnight
	prettyClassName = "".join(word.title() for i, word in enumerate(className.split("-")))
	outfile.write("message {}Talents {{\n".format(prettyClassName))

	talentIdx = 1
	for treeIdx, tree in enumerate(trees):
		title = tree.find_element(By.XPATH, "./div/b").text
		outfile.write("\t// {}\n".format(title))

		treeTalentsData = []
		talents = tree.find_elements(By.CLASS_NAME, "ctc-tree-talent")
		for talent in talents:
			max_points = int(talent.get_attribute("data-max-points"))
			field_type = "bool" if max_points == 1 else "int32"

			link = talent.find_element(By.XPATH, "./div/a").get_attribute("href")
			name = "_".join(word for i, word in enumerate(link.split("/")[-1].split("-")))

			print("Talent: " + name)
			treeTalentsData.append({
				"row": int(talent.get_attribute("data-row")),
				"col": int(talent.get_attribute("data-col")),
				"name": name,
				"field_type": field_type,
			})

		treeTalentsData.sort(key=lambda data: data["row"] * 4 + data["col"])
		for data in treeTalentsData:
			outfile.write("\t{} {} = {};\n".format(data["field_type"], data["name"], talentIdx))
			talentIdx += 1

		if treeIdx != len(trees) - 1:
			outfile.write("\n")

	outfile.write("}")
