package screenconfig

import "errors"

// ErrActionsMixedWithAddedRemoved es el error de validacion cuando una
// screen_instance declara simultaneamente la lista legacy "actions" y
// los campos compositivos "actions_added"/"actions_removed". El composer
// no acepta el hibrido: el seed debe optar por uno u otro modelo.
var ErrActionsMixedWithAddedRemoved = errors.New("screen_instance: slot_data.actions cannot coexist with actions_added/actions_removed")
