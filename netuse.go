package netresource

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var procWNetAddConnection2W = mpr.NewProc("WNetAddConnection2W")
var procWNetCancelConnection = mpr.NewProc("WNetCancelConnection2W")

// WNetAddConnection2 makes a connection to a network resource
// and can redirect a local device to the network resource
//     remote - UNC-Path like `\\localhost\C$`
//     local - local-drive like `X:`
//     user - username. When it is "", default username is used.
//     pass - password. When it is "", default passwors is used.
func WNetAddConnection2(remote, local, user, pass string) (err error) {
	var rs NetResource

	rs.localName, err = windows.UTF16PtrFromString(local)
	if err != nil {
		return
	}
	rs.remoteName, err = windows.UTF16PtrFromString(remote)
	if err != nil {
		return
	}
	var _user *uint16
	if user == "" {
		_user = nil
	} else {
		_user, err = windows.UTF16PtrFromString(user)
		if err != nil {
			return
		}
	}
	var _pass *uint16
	if pass == "" {
		_pass = nil
	} else {
		_pass, err = windows.UTF16PtrFromString(pass)
		if err != nil {
			return
		}
	}

	rc, _, err := procWNetAddConnection2W.Call(
		uintptr(unsafe.Pointer(&rs)),
		uintptr(unsafe.Pointer(_pass)),
		uintptr(unsafe.Pointer(_user)),
		0)

	if rc != 0 {
		return err
	}
	return nil
}

// WNetCancelConnection2 cancels an existing network connection.
//    update - true: updates connection as not a persistent one
//    force - true: disconnect even if open process exists.
func WNetCancelConnection2(name string, update bool, force bool) error {
	_name, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	var _update uintptr
	if update {
		_update = CONNECT_UPDATE_PROFILE
	}
	var _force uintptr
	if force {
		_force = 1
	}
	rc, _, err := procWNetCancelConnection.Call(
		uintptr(unsafe.Pointer(_name)),
		_update,
		_force)
	if rc != 0 {
		return err
	}
	return nil
}
