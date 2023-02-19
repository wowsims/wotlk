import { TargetInput } from "./proto/common";

export class TargetInputs {
	private readonly targetInputs: Array<TargetInput>;

	constructor(targetInputs?: Array<TargetInput>) {
		this.targetInputs = TargetInputs.initTargetInputsArray(targetInputs);
	}

	private static initTargetInputsArray(newTargetInputs?: Array<TargetInput>): Array<TargetInput> {
		return newTargetInputs?.slice(0, newTargetInputs.length) || [];
	}

	private static targetInputEqual(lhs: TargetInput, rhs: TargetInput): boolean {
		return lhs?.label == rhs?.label && lhs?.inputType == rhs?.inputType && lhs?.boolValue == rhs?.boolValue && lhs?.numberValue == rhs?.numberValue;
	}

	equals(other: TargetInputs): boolean {
		if (this.targetInputs?.length != other.targetInputs?.length) {
			return false;
		}
		return this.targetInputs.every((newTargetInput, inputIdx) => TargetInputs.targetInputEqual(newTargetInput, other.getTargetInput(inputIdx)));
	}

	getLength(): number {
		return this.targetInputs?.length;
	}

	getTargetInput(index: number): TargetInput {
		return this.targetInputs[index];
	}

	asArray(): Array<TargetInput> {
		return this.targetInputs.slice();
	}

	hasInputs(): boolean {
		return this.targetInputs?.length > 0;
	}
}