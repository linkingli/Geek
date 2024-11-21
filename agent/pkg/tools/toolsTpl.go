package tools

const AddToolName = `AddTool`

const AddToolDescription = `
Use this tool for addition calculations.
	example:
		1+2 =?
	then Action Input is: 1,2
`

const AddToolParam = `{"type":"object","properties":{"numbers":{"type":"array","items":{"type":"integer"}}}}`

const SubToolName = `SubTool`

const SubToolDescription = `
Use this tool for subtraction calculations.
	example:
		1-2 =?
	then Action Input is: 1,2
`

const SubToolParam = `{"type":"object","properties":{"numbers":{"type":"array","items":{"type":"integer"}}}}`
