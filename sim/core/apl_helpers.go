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

func (rot *APLRotation) getSourceUnit(ref *proto.UnitReference) UnitReference {
	if ref == nil || ref.Type == proto.UnitReference_Unknown {
		return NewUnitReference(&proto.UnitReference{Type: proto.UnitReference_Self}, rot.unit)
	} else {
		unitRef := NewUnitReference(ref, rot.unit)
		if unitRef.Get() == nil {
			rot.validationWarning("No unit found matching reference: %s", ref)
		}
		return unitRef
	}
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

func (rot *APLRotation) aplGetAura(sourceRef *proto.UnitReference, auraId *proto.ActionID) AuraReference {
	sourceUnit := rot.getSourceUnit(sourceRef)
	if sourceUnit.Get() == nil {
		return AuraReference{}
	}

	aura := NewAuraReference(sourceUnit, auraId)
	if aura.Get() == nil {
		rot.validationWarning("No aura found on %s for: %s", sourceUnit.Get().Label, ProtoToActionID(auraId))
	}
	return aura
}

func (rot *APLRotation) aplGetProcAura(sourceRef *proto.UnitReference, auraId *proto.ActionID) AuraReference {
	sourceUnit := rot.getSourceUnit(sourceRef)
	if sourceUnit.Get() == nil {
		return AuraReference{}
	}

	aura := NewIcdAuraReference(sourceUnit, auraId)
	if aura.Get() == nil {
		rot.validationWarning("No aura found on %s for: %s", sourceUnit.Get().Label, ProtoToActionID(auraId))
	}
	return aura
}

func (rot *APLRotation) aplGetSpell(spellId *proto.ActionID) *Spell {
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
		spell = rot.unit.GetSpell(actionID)
	}

	if spell == nil {
		rot.validationWarning("%s does not know spell %s", rot.unit.Label, actionID)
	}
	return spell
}

func (rot *APLRotation) aplGetDot(spellId *proto.ActionID) *Dot {
	spell := rot.aplGetSpell(spellId)

	if spell == nil {
		return nil
	} else if spell.AOEDot() != nil {
		return spell.AOEDot()
	} else {
		return spell.CurDot()
	}
}

func (rot *APLRotation) aplGetMultidotSpell(spellId *proto.ActionID) *Spell {
	spell := rot.aplGetSpell(spellId)
	if spell == nil {
		return nil
	} else if spell.CurDot() == nil {
		rot.validationWarning("Spell %s does not have an associated DoT", ProtoToActionID(spellId))
		return nil
	}
	return spell
}
