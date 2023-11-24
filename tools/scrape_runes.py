#!/usr/bin/python

# This tool generates the classic SoD runes data

import sys
import requests
import math

from typing import List

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from webdriver_manager.chrome import ChromeDriverManager

if len(sys.argv) < 2:
	raise Exception("Missing arguments, expected output_file_path")

output_file_path = sys.argv[1]

driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()))
wait = WebDriverWait(driver, 10)
element_locator = (By.ID, "data-tree-switcher")

def _get_spell_id_from_link(link):
	return int(link.split("/")[-2].split("=")[-1])


def get_engraving_ids() -> List[int]:
	driver.get(f"https://www.wowhead.com/classic/search?q=engrave")
	wait.until(EC.presence_of_element_located(element_locator))

	abilities = driver.find_elements(By.ID, "tab-abilities")

	if len(abilities) == 0:
		print(f"Engravings missing ability tab.")
		return []

	abilities = abilities[0]
	pages = int(abilities.find_element(By.CLASS_NAME, "listview-nav").find_element(By.CSS_SELECTOR, 'b:last-child').text)/50
	pages = math.ceil(pages)
	all_ids = []

	for page in range(pages):
		print(f'Loading page {page} for runes...')
		driver.get(f"https://www.wowhead.com/classic/search?q=engrave#abilities:{page*50}")
		driver.refresh()
		wait.until(EC.presence_of_element_located(element_locator))
		abilities_tab = driver.find_element(By.ID, "tab-abilities")
		rows = abilities_tab.find_elements(By.CLASS_NAME, "listview-row")
		all_ids.extend([_get_spell_id_from_link(row.find_element(By.CLASS_NAME, "listview-cleartext").get_attribute("href"))
			for row in rows])

	driver.quit()
	return all_ids

# id,tooltip_json
to_export = []

engraving_ids = get_engraving_ids()

print(f"Export Count ({len(engraving_ids)}) {engraving_ids}")

for id in engraving_ids:
	url = f"https://nether.wowhead.com/classic/tooltip/spell/{id}"
	result = requests.get(url)

	if result.status_code == 200:
		response_json = result.text
		to_export.append([id, response_json])
	else:
		print(f"Request for id {id} failed with status code: {result.status_code}")

output_string = '\n'.join([str(','.join([str(inner_elem) for inner_elem in elem])) for elem in to_export])

with open(output_file_path, "w") as outfile:
 	outfile.write(output_string)
