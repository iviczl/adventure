package constants

import (
	"path/filepath"
	"text-adventure/session"
	"text-adventure/types"

	"gorm.io/gorm"
)

const (
	CHANGE_PLAYER_ATTRIBUTES               types.ActionOperation = "cplat"
	CHANGE_PLAYER_ABILITIES                types.ActionOperation = "cplab"
	ADD_PLAYER_ABILITIES                   types.ActionOperation = "aplab"
	ERASE_PLAYER_ABILITIES                 types.ActionOperation = "eplab"
	CHANGE_POSITION                        types.ActionOperation = "cp"
	CHANGE_POSITION_DESCRIPTION            types.ActionOperation = "cpd"
	APPEND_POSITION_TEMPORARY_DESCRIPTION  types.ActionOperation = "aptd"
	CHANGE_POSITION_TEMPORARY_DESCRIPTION  types.ActionOperation = "cptd"
	PREPEND_POSITION_TEMPORARY_DESCRIPTION types.ActionOperation = "pptd"
	CHANGE_POSITION_VISITED                types.ActionOperation = "cpv"
	CHANGE_ACTION_ACTIVE                   types.ActionOperation = "caa"
	CHANGE_ACTION_VISIBLE                  types.ActionOperation = "cav"
	CHANGE_ITEM_DESCRIPTION                types.ActionOperation = "cad"
	CHANGE_ITEM_STATE                      types.ActionOperation = "cis"
	CHANGE_NPC_STATE                       types.ActionOperation = "cns"
	CHANGE_NPC_INTERACTED                  types.ActionOperation = "cni"
	CHANGE_NPC_ATTRIBUTES                  types.ActionOperation = "cnat"
	CHANGE_NPC_ABILITIES                   types.ActionOperation = "cnab"
	ATTACK_NPC                             types.ActionOperation = "an"
	PICK_UP_ITEM                           types.ActionOperation = "pui"
	PUT_DOWN_ITEM                          types.ActionOperation = "pdi"
	MOVE_ITEM                              types.ActionOperation = "moi"
	ERASE_ITEM                             types.ActionOperation = "eri"
	CONDITIONAL                            types.ActionOperation = "con"
	SWITCH_CONDITIONAL                     types.ActionOperation = "scon"
	LIST                                   types.ActionOperation = "ls"
	RANDOM                                 types.ActionOperation = "ran"
)

var SessionStore *session.InMemoryStore
var DbClient *gorm.DB

var DbPath string = filepath.Join(".", "db", "text-adventure.db")

const SecretKey = "my_secret_key"
const SessionName = "session"
const SessionSecret = "WYVOEWXHLUDB34DAELL2LDYNJSLEZT5WZRLPLHIQE5JXOGYCGIZQ"
