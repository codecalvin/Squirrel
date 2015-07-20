package controllers

type StatusCode int

const (
	START_CODE StatusCode = iota
	OK = 200
	BadRequest = 400
	Unauth = iota
)