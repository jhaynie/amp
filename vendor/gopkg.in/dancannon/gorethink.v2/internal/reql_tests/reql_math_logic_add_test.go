// Code generated by gen_tests.py and process_polyglot.py.
// Do not edit this file directly.
// The template for this file is located at:
// ../template.go.tpl
package reql_tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	r "gopkg.in/gorethink/gorethink.v2"
	"gopkg.in/gorethink/gorethink.v2/internal/compare"
)

// Tests for basic usage of the add operation
func TestMathLogicAddSuite(t *testing.T) {
	suite.Run(t, new(MathLogicAddSuite))
}

type MathLogicAddSuite struct {
	suite.Suite

	session *r.Session
}

func (suite *MathLogicAddSuite) SetupTest() {
	suite.T().Log("Setting up MathLogicAddSuite")
	// Use imports to prevent errors
	_ = time.Time{}
	_ = compare.AnythingIsFine

	session, err := r.Connect(r.ConnectOpts{
		Address: url,
	})
	suite.Require().NoError(err, "Error returned when connecting to server")
	suite.session = session

	r.DBDrop("test").Exec(suite.session)
	err = r.DBCreate("test").Exec(suite.session)
	suite.Require().NoError(err)
	err = r.DB("test").Wait().Exec(suite.session)
	suite.Require().NoError(err)

}

func (suite *MathLogicAddSuite) TearDownSuite() {
	suite.T().Log("Tearing down MathLogicAddSuite")

	if suite.session != nil {
		r.DB("rethinkdb").Table("_debug_scratch").Delete().Exec(suite.session)
		r.DBDrop("test").Exec(suite.session)

		suite.session.Close()
	}
}

func (suite *MathLogicAddSuite) TestCases() {
	suite.T().Log("Running MathLogicAddSuite: Tests for basic usage of the add operation")

	{
		// math_logic/add.yaml line #3
		/* 2 */
		var expected_ int = 2
		/* r.add(1, 1) */

		suite.T().Log("About to run line #3: r.Add(1, 1)")

		runAndAssert(suite.Suite, expected_, r.Add(1, 1), suite.session, r.RunOpts{
			GeometryFormat: "raw",
			GroupFormat:    "map",
		})
		suite.T().Log("Finished running line #3")
	}

	{
		// math_logic/add.yaml line #8
		/* 2 */
		var expected_ int = 2
		/* r.expr(1) + 1 */

		suite.T().Log("About to run line #8: r.Expr(1).Add(1)")

		runAndAssert(suite.Suite, expected_, r.Expr(1).Add(1), suite.session, r.RunOpts{
			GeometryFormat: "raw",
			GroupFormat:    "map",
		})
		suite.T().Log("Finished running line #8")
	}

	{
		// math_logic/add.yaml line #9
		/* 2 */
		var expected_ int = 2
		/* 1 + r.expr(1) */

		suite.T().Log("About to run line #9: r.Add(1, r.Expr(1))")

		runAndAssert(suite.Suite, expected_, r.Add(1, r.Expr(1)), suite.session, r.RunOpts{
			GeometryFormat: "raw",
			GroupFormat:    "map",
		})
		suite.T().Log("Finished running line #9")
	}

	{
		// math_logic/add.yaml line #10
		/* 2 */
		var expected_ int = 2
		/* r.expr(1).add(1) */

		suite.T().Log("About to run line #10: r.Expr(1).Add(1)")

		runAndAssert(suite.Suite, expected_, r.Expr(1).Add(1), suite.session, r.RunOpts{
			GeometryFormat: "raw",
			GroupFormat:    "map",
		})
		suite.T().Log("Finished running line #10")
	}

	{
		// math_logic/add.yaml line #16
		/* 0 */
		var expected_ int = 0
		/* r.expr(-1) + 1 */

		suite.T().Log("About to run line #16: r.Expr(-1).Add(1)")

		runAndAssert(suite.Suite, expected_, r.Expr(-1).Add(1), suite.session, r.RunOpts{
			GeometryFormat: "raw",
			GroupFormat:    "map",
		})
		suite.T().Log("Finished running line #16")
	}

	{
		// math_logic/add.yaml line #21
		/* 10.25 */
		var expected_ float64 = 10.25
		/* r.expr(1.75) + 8.5 */

		suite.T().Log("About to run line #21: r.Expr(1.75).Add(8.5)")

		runAndAssert(suite.Suite, expected_, r.Expr(1.75).Add(8.5), suite.session, r.RunOpts{
			GeometryFormat: "raw",
			GroupFormat:    "map",
		})
		suite.T().Log("Finished running line #21")
	}

	{
		// math_logic/add.yaml line #27
		/* '' */
		var expected_ string = ""
		/* r.expr('') + '' */

		suite.T().Log("About to run line #27: r.Expr('').Add('')")

		runAndAssert(suite.Suite, expected_, r.Expr("").Add(""), suite.session, r.RunOpts{
			GeometryFormat: "raw",
			GroupFormat:    "map",
		})
		suite.T().Log("Finished running line #27")
	}

	{
		// math_logic/add.yaml line #32
		/* 'abcdef' */
		var expected_ string = "abcdef"
		/* r.expr('abc') + 'def' */

		suite.T().Log("About to run line #32: r.Expr('abc').Add('def')")

		runAndAssert(suite.Suite, expected_, r.Expr("abc").Add("def"), suite.session, r.RunOpts{
			GeometryFormat: "raw",
			GroupFormat:    "map",
		})
		suite.T().Log("Finished running line #32")
	}

	{
		// math_logic/add.yaml line #52
		/* err("ReqlQueryLogicError", "Expected type NUMBER but found STRING.", [1]) */
		var expected_ Err = err("ReqlQueryLogicError", "Expected type NUMBER but found STRING.")
		/* r.expr(1) + 'a' */

		suite.T().Log("About to run line #52: r.Expr(1).Add('a')")

		runAndAssert(suite.Suite, expected_, r.Expr(1).Add("a"), suite.session, r.RunOpts{
			GeometryFormat: "raw",
			GroupFormat:    "map",
		})
		suite.T().Log("Finished running line #52")
	}

	{
		// math_logic/add.yaml line #57
		/* err("ReqlQueryLogicError", "Expected type STRING but found NUMBER.", [1]) */
		var expected_ Err = err("ReqlQueryLogicError", "Expected type STRING but found NUMBER.")
		/* r.expr('a') + 1 */

		suite.T().Log("About to run line #57: r.Expr('a').Add(1)")

		runAndAssert(suite.Suite, expected_, r.Expr("a").Add(1), suite.session, r.RunOpts{
			GeometryFormat: "raw",
			GroupFormat:    "map",
		})
		suite.T().Log("Finished running line #57")
	}

	{
		// math_logic/add.yaml line #62
		/* err("ReqlQueryLogicError", "Expected type ARRAY but found NUMBER.", [1]) */
		var expected_ Err = err("ReqlQueryLogicError", "Expected type ARRAY but found NUMBER.")
		/* r.expr([]) + 1 */

		suite.T().Log("About to run line #62: r.Expr([]interface{}{}).Add(1)")

		runAndAssert(suite.Suite, expected_, r.Expr([]interface{}{}).Add(1), suite.session, r.RunOpts{
			GeometryFormat: "raw",
			GroupFormat:    "map",
		})
		suite.T().Log("Finished running line #62")
	}
}
