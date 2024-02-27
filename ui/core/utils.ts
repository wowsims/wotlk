// Returns if the two items are equal, or if both are null / undefined.
export function equalsOrBothNull<T>(a: T, b: T, comparator?: (_a: NonNullable<T>, _b: NonNullable<T>) => boolean): boolean {
	if (a == null && b == null)
		return true;

	if (a == null || b == null)
		return false;

	return (comparator || ((_a: NonNullable<T>, _b: NonNullable<T>) => a == b))(a!, b!);
}

// Default comparator function for strings. Used with functions like Array.sort().
export function stringComparator(a: string, b: string): number {
	if (a < b) {
		return -1;
	} else if (b < a) {
		return 1;
	} else {
		return 0;
	}
}

// Sorts an objectArray by a property. Returns a new array.
// Can be called recursively.
export function sortByProperty(objArray: any[], prop: string) {
	if (!Array.isArray(objArray)) throw new Error('FIRST ARGUMENT NOT AN ARRAY');
	const clone = objArray.slice(0);
	const direct = arguments.length > 2 ? arguments[2] : 1; //Default to ascending
	const propPath = (prop.constructor === Array) ? prop : prop.split('.');
	clone.sort(function(a, b) {
		for (const p in propPath) {
			if (a[propPath[p]] && b[propPath[p]]) {
				a = a[propPath[p]];
				b = b[propPath[p]];
			}
		}
		// convert numeric strings to integers
		a = a.toString().match(/^\d+$/) ? +a : a;
		b = b.toString().match(/^\d+$/) ? +b : b;
		return ((a < b) ? -1 * direct : ((a > b) ? 1 * direct : 0));
	});
	return clone;
}

export function sum(arr: Array<number>): number {
	return arr.reduce((total, cur) => total + cur, 0);
}

// Returns the index of maximum value, or null if empty.
export function maxIndex(arr: Array<number>): number | null {
	return arr.reduce((cur, v, i, arr) => v > arr[cur] ? i : cur, 0);
}

// Swaps two elements in the given array.
export function swap<T>(arr: Array<T>, i: number, j: number) {
	[arr[i], arr[j]] = [arr[j], arr[i]];
}

// Returns a new array containing only elements present in both a and b.
export function arrayEquals<T>(a: Array<T>, b: Array<T>, comparator?: (a: T, b: T) => boolean): boolean {
	comparator = comparator || ((a: T, b: T) => a == b);
	return a.length == b.length && a.every((val, i) => comparator!(val, b[i]));
}

// Returns a new array containing only elements present in both a and b.
export function intersection<T>(a: Array<T>, b: Array<T>): Array<T> {
	return a.filter(value => b.includes(value));
}

// Returns a new array containing only distinct elements of arr.
// comparator should return true if the two elements are considered equal, and false otherwise.
export function distinct<T>(arr: Array<T>, comparator?: (a: T, b: T) => boolean): Array<T> {
	comparator = comparator || ((a: T, b: T) => a == b);
	const distinctArr: Array<T> = [];
	arr.forEach(val => {
		if (distinctArr.find(dVal => comparator!(dVal, val)) == null) {
			distinctArr.push(val);
		}
	});
	return distinctArr;
}

// Splits an array into buckets, where elements are placed in the same bucket if the
// toString function returns the same value.
export function bucket<T>(arr: Array<T>, toString: (val: T) => string): Record<string, Array<T>> {
	const buckets: Record<string, Array<T>> = {};
	arr.forEach(val => {
		const valString = toString(val);
		if (buckets[valString]) {
			buckets[valString].push(val);
		} else {
			buckets[valString] = [val];
		}
	});
	return buckets;
}

export function stDevToConf90(stDev: number, N: number) {
	return 1.645 * stDev / Math.sqrt(N);
}

export async function wait(ms: number): Promise<void> {
	return new Promise(resolve => setTimeout(resolve, ms));
}

// Only works for numeric enums
export function getEnumValues<E>(enumType: any): Array<E> {
	return Object.keys(enumType)
		.filter(key => !isNaN(Number(enumType[key])))
		.map(key => parseInt(enumType[key]) as unknown as E);
}

// Whether a click event was a right click.
export function isRightClick(event: MouseEvent): boolean {
	return event.button == 2;
}

// Converts from '#ffffff' --> 'rgba(255, 255, 255, alpha)'
export function hexToRgba(hex: string, alpha: number): string {
	if (/^#([A-Fa-f0-9]{3}){1,2}$/.test(hex)) {
		let parts = hex.substring(1).split('');
		if (parts.length == 3) {
			parts = [parts[0], parts[0], parts[1], parts[1], parts[2], parts[2]];
		}
		const c: any = '0x' + parts.join('');
		return 'rgba(' + [(c >> 16) & 255, (c >> 8) & 255, c & 255].join(',') + ',' + alpha + ')';
	}
	throw new Error('Invalid hex color: ' + hex);
}

export function camelToSnakeCase(str: string): string {
	let result = str.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`);
	if (result.startsWith('_')) {
		result = result.substring(1);
	}
	return result;
}

export function downloadJson(json: any, fileName: string) {
	downloadString(JSON.stringify(json, null, 2), fileName);
}
export function downloadString(data: string, fileName: string) {
	const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(data);
	const downloadAnchorNode = document.createElement('a');
	downloadAnchorNode.setAttribute("href", dataStr);
	downloadAnchorNode.setAttribute("download", fileName);
	document.body.appendChild(downloadAnchorNode); // required for firefox
	downloadAnchorNode.click();
	downloadAnchorNode.remove();
}

export function formatDeltaTextElem(elem: HTMLElement, before: number, after: number, precision: number, lowerIsBetter?: boolean, noColor?: boolean) {
	const delta = after - before;
	let deltaStr = delta.toFixed(precision);
	if (delta >= 0) {
		deltaStr = '+' + deltaStr;
	}
	elem.textContent = deltaStr;

	if (noColor || delta == 0) {
		elem.classList.remove('positive');
		elem.classList.remove('negative');
	} else if (delta > 0 != Boolean(lowerIsBetter)) {
		elem.classList.remove('negative');
		elem.classList.add('positive');
	} else {
		elem.classList.remove('positive');
		elem.classList.add('negative');
	}
}

// Returns all N pick K permutations of the elements in arr of size N.
export function permutations<T>(arr: Array<T>, k: number): Array<Array<T>> {
	if (k == 0) {
		return [];
	} else if (k == 1) {
		return arr.map(v => [v]);
	} else {
		return arr.map((v, i) => {
			const withoutThisElem = arr.slice();
			withoutThisElem.splice(i, 1);
			const permutationsWithoutThisElem = permutations(withoutThisElem, k - 1);
			return permutationsWithoutThisElem.map(perm => [v].concat(perm));
		}).flat();
	}
}

// Returns all N choose K combinations of the elements in arr of size N.
export function combinations<T>(arr: Array<T>, k: number, comparator?: (_a: T, _b: T) => number): Array<Array<T>> {
	const perms = permutations(arr, k);
	const sorted = perms.map(permutation => permutation.sort(comparator));

	const equals: ((_a: T, _b: T) => boolean) = comparator ? ((a, b) => comparator(a, b) == 0) : ((a, b) => a == b);
	return distinct(sorted, (permutationA, permutationB) => permutationA.every((elem, i) => equals(elem, permutationB[i])));
}

// Returns all N pick K permutations of the elements in arr of size N, allowing duplicates.
export function permutationsWithDups<T>(arr: Array<T>, k: number): Array<Array<T>> {
	if (k == 0) {
		return [];
	} else if (k == 1) {
		return arr.map(v => [v]);
	} else {
		const smaller = permutationsWithDups(arr, k - 1);
		return arr.map(v => {
			return smaller.map(permutation => {
				const newPerm = permutation.slice();
				newPerm.push(v);
				return newPerm;
			});
		}).flat();
	}
}

// Returns all N choose K combinations of the elements in arr of size N, allowing duplicates.
export function combinationsWithDups<T>(arr: Array<T>, k: number): Array<Array<T>> {
	const perms = permutationsWithDups(arr, k);
	const sorted = perms.map(permutation => permutation.sort());
	return distinct(sorted, (permutationA, permutationB) => permutationA.every((elem, i) => elem == permutationB[i]));
}

// Converts a Uint8Array into a hex string.
export function buf2hex(data: Uint8Array): string {
	return [...data]
		.map(x => x.toString(16).padStart(2, '0'))
		.join('');
}

const randomStringChars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+';
export function randomString(len?: number): string {
	let str = '';
	const strLen = len || 10;
	for (let i = 0; i < strLen; i++) {
		str += randomStringChars[Math.floor(Math.random() * randomStringChars.length)];
	}
	return str;
}

// Allows replacement of stringified objects based on the key and path.
// If handler returns a string, that string is used. Otherwise, the normal JSON.stringify result is returned.
export function jsonStringifyCustom(value: any, indent: number, handler: (value: any, path: Array<string>) => string | undefined | void): string {
	const indentStr = ' '.repeat(indent);
	return jsonStringifyCustomHelper(value, indentStr, [], handler);
}
function jsonStringifyCustomHelper(value: any, indentStr: string, path: Array<string>, handler: (value: any, path: Array<string>) => string | undefined | void): string {
	const handlerResult = handler(value, path);
	if (handlerResult != null) {
		return handlerResult;
	}

	if (!(value instanceof Object)) {
		return JSON.stringify(value);
	} else if (value instanceof Array) {
		let str = '[\n';
		const lines = value.map((e, i) => `${indentStr.repeat(path.length + 1)}${jsonStringifyCustomHelper(e, indentStr, path.slice().concat([i + '']), handler)}${i == value.length - 1 ? '' : ','}\n`);
		str += lines.join('');
		str += indentStr.repeat(path.length) + ']';
		return str;
	} else { // Object
		let str = '{\n';
		const len = Object.keys(value).length;
		const lines = Object.entries(value).map(([fieldKey, fieldValue], i) => `${indentStr.repeat(path.length + 1)}"${fieldKey}": ${jsonStringifyCustomHelper(fieldValue, indentStr, path.slice().concat([fieldKey]), handler)}${i == len - 1 ? '' : ','}\n`);
		str += lines.join('');
		str += indentStr.repeat(path.length) + '}';
		return str;
	}
}

// Pretty-prints the value in JSON form, but does not prettify (flattens) sub-values where handler returns true.
export function jsonStringifyWithFlattenedPaths(value: any, indent: number, handler: (value: any, path: Array<string>) => boolean): string {
	return jsonStringifyCustom(value, indent, (value, path) => handler(value, path) ? JSON.stringify(value) : undefined);
}

export function htmlDecode(input: string) {
	const doc = new DOMParser().parseFromString(input, "text/html");
	return doc.documentElement.textContent;
}
