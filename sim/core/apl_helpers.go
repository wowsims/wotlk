package core

import (
	"github.com/wowsims/wotlk/sim/core/proto"
)

// Struct for handling unit references, to account for values that can
// change dynamically (e.g. CurrentTarget).
type UnitReference struct {
	fixedUnit       *Unit
	curTargetSource *Unit
}

func (ur UnitReference) Get() *Unit {
	if ur.fixedUnit != nil {
		return ur.fixedUnit
	} else if ur.curTargetSource != nil {
		return ur.curTargetSource.CurrentTarget
	} else {
		return nil
	}
}

func (ur *UnitReference) String() string {
	return ur.Get().Label
}

func NewUnitReference(ref *proto.UnitReference, contextUnit *Unit) UnitReference {
	if ref == nil || ref.Type == proto.UnitReference_Unknown {
		return UnitReference{}
	} else if ref.Type == proto.UnitReference_CurrentTarget {
		return UnitReference{
			curTargetSource: contextUnit,
		}
	} else {
		return UnitReference{
			fixedUnit: contextUnit.GetUnit(ref),
		}
	}
}

func (rot *APLRotation) getUnit(ref *proto.UnitReference, defaultRef *proto.UnitReference) UnitReference {
	if ref == nil || ref.Type == proto.UnitReference_Unknown {
		return NewUnitReference(defaultRef, rot.unit)
	} else {
		unitRef := NewUnitReference(ref, rot.unit)
		if unitRef.Get() == nil {
			rot.ValidationWarning("No unit found matching reference: %s", ref)
		}
		return unitRef
	}
}
func (rot *APLRotation) GetSourceUnit(ref *proto.UnitReference) UnitReference {
	return rot.getUnit(ref, &proto.UnitReference{Type: proto.UnitReference_Self})
}
func (rot *APLRotation) GetTargetUnit(ref *proto.UnitReference) UnitReference {
	return rot.getUnit(ref, &proto.UnitReference{Type: proto.UnitReference_CurrentTarget})
}

type AuraReference struct {
	fixedAura *Aura

	curTargetSource *Unit
	curTargetAuras  AuraArray
}

func (ar *AuraReference) Get() *Aura {
	if ar.fixedAura != nil {
		return ar.fixedAura
	} else if ar.curTargetSource != nil {
		return ar.curTargetAuras.Get(ar.curTargetSource.CurrentTarget)
	} else {
		return nil
	}
}

func (ar *AuraReference) String() string {
	return ar.Get().ActionID.String()
}

func newAuraReferenceHelper(sourceUnit UnitReference, auraId *proto.ActionID, auraGetter func(*Unit, ActionID) *Aura) AuraReference {
	if sourceUnit.Get() == nil {
		return AuraReference{}
	} else if sourceUnit.fixedUnit != nil {
		return AuraReference{
			fixedAura: auraGetter(sourceUnit.fixedUnit, ProtoToActionID(auraId)),
		}
	} else {
		auras := make([]*Aura, len(sourceUnit.Get().Env.AllUnits))
		for _, unit := range sourceUnit.Get().Env.AllUnits {
			auras[unit.UnitIndex] = auraGetter(unit, ProtoToActionID(auraId))
		}
		return AuraReference{
			curTargetSource: sourceUnit.curTargetSource,
			curTargetAuras:  auras,
		}
	}
}
func NewAuraReference(sourceUnit UnitReference, auraId *proto.ActionID) AuraReference {
	return newAuraReferenceHelper(sourceUnit, auraId, func(unit *Unit, actionID ActionID) *Aura { return unit.GetAuraByID(actionID) })
}
func NewIcdAuraReference(sourceUnit UnitReference, auraId *proto.ActionID) AuraReference {
	return newAuraReferenceHelper(sourceUnit, auraId, func(unit *Unit, actionID ActionID) *Aura { return unit.GetIcdAuraByID(actionID) })
}

func (rot *APLRotation) GetAPLAura(sourceUnit UnitReference, auraId *proto.ActionID) AuraReference {
	if sourceUnit.Get() == nil {
		return AuraReference{}
	}

	aura := NewAuraReference(sourceUnit, auraId)
	if aura.Get() == nil {
		rot.ValidationWarning("No aura found on %s for: %s", sourceUnit.Get().Label, ProtoToActionID(auraId))
	}
	return aura
}

func (rot *APLRotation) GetAPLICDAura(sourceUnit UnitReference, auraId *proto.ActionID) AuraReference {
	if sourceUnit.Get() == nil {
		return AuraReference{}
	}

	aura := NewIcdAuraReference(sourceUnit, auraId)
	if aura.Get() == nil {
		rot.ValidationWarning("No aura found on %s for: %s", sourceUnit.Get().Label, ProtoToActionID(auraId))
	}
	return aura
}

func (rot *APLRotation) GetAPLSpell(spellId *proto.ActionID) *Spell {
	actionID := ProtoToActionID(spellId)
	var spell *Spell

	if actionID.IsOtherAction(proto.OtherAction_OtherActionPotion) {
		if rot.parsingPrepull {
			for _, s := range rot.unit.Spellbook {
				if s.Flags.Matches(SpellFlagPrepullPotion) {
					spell = s
					break
				}
			}
		} else {
			for _, s := range rot.unit.Spellbook {
				if s.Flags.Matches(SpellFlagCombatPotion) {
					spell = s
					break
				}
			}
		}
	} else {
		// Prefer spells marked with APL, but fallback to unmarked spells.
		var aplSpell *Spell
		for _, s := range rot.unit.Spellbook {
			if s.ActionID.SameAction(actionID) && s.Flags.Matches(SpellFlagAPL) {
				aplSpell = s
				break
			}
		}
		if aplSpell == nil {
			spell = rot.unit.GetSpell(actionID)
		} else {
			spell = aplSpell
		}
	}

	if spell == nil {
		rot.ValidationWarning("%s does not know spell %s", rot.unit.Label, actionID)
	}
	return spell
}

func (rot *APLRotation) GetTargetAPLSpell(spellId *proto.ActionID, targetUnit UnitReference) *Spell {
	actionID := ProtoToActionID(spellId)
	target := targetUnit.Get()
	spell := target.GetSpell(actionID)

	if spell == nil {
		rot.ValidationWarning("%s does not know spell %s", target.Label, actionID)
	}
	return spell
}

func (rot *APLRotation) GetAPLDot(targetUnit UnitReference, spellId *proto.ActionID) *Dot {
	spell := rot.GetAPLSpell(spellId)

	if spell == nil {
		return nil
	} else if spell.AOEDot() != nil {
		return spell.AOEDot()
	} else {
		target := targetUnit.Get()
		if target != nil {
			return spell.Dot(target)
		} else {
			return spell.CurDot()
		}
	}
}

func (rot *APLRotation) GetAPLMultidotSpell(spellId *proto.ActionID) *Spell {
	spell := rot.GetAPLSpell(spellId)
	if spell == nil {
		return nil
	} else if spell.CurDot() == nil {
		rot.ValidationWarning("Spell %s does not have an associated DoT", ProtoToActionID(spellId))
		return nil
	}
	return spell
}

func (rot *APLRotation) GetAPLMultishieldSpell(spellId *proto.ActionID) *Spell {
	spell := rot.GetAPLSpell(spellId)
	if spell == nil {
		return nil
	} else if spell.Shield(spell.Unit) == nil {
		rot.ValidationWarning("Spell %s does not have an associated Shield", ProtoToActionID(spellId))
		return nil
	}
	return spell
}
