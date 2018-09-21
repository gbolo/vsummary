// Package datatablessrv handles the server side processing of an AJAX request for DataTables
// For details on the parameters and the results, read the datatables documentation at
// https://datatables.net/manual/server-side
package server

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	//_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Copyright (c) 2017 Escape Velocity, Inc.
// Copyright (c) 2018 gbolo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

// ErrNotDataTablesReq indicates that this is not being requested by Datatables
var ErrNotDataTablesReq = errors.New("Not a DataTables request")

// SortDir is the direction of the sort (ascending/descending)
type SortDir int

// Database Instance
var Dbx sqlx.DB

const (
	// Asc for ascending sorting
	Asc SortDir = iota
	// Desc for descending sorting
	Desc
)

// OrderInfo tracks the list of columns to sort by and in which direction to sort them.
type OrderInfo struct {
	// ColNum indicates Which column to apply sorting to (zero based index to the Columns data)
	ColNum int
	// Direction tells us which way to sort
	Direction SortDir
}

// ColData tracks all of the columns requested by DataTables
type ColData struct {
	// columns[i][name] Column's name, as defined by columns.name.
	Name string
	// columns[i][data] Column's data source, as defined by columns.data.
	// It is poss
	Data string
	// columns[i][searchable]	boolean	Flag to indicate if this column is searchable (true) or not (false).
	// This is controlled by columns.searchable.
	Searchable bool
	// columns[i][orderable] Flag to indicate if this column is orderable (true) or not (false).
	// This is controlled by columns.orderable.
	Orderable bool
	// columns[i][search][value] Search value to apply to this specific column.
	Searchval string
	// columns[i][search][regex]
	// Flag to indicate if the search term for this column should be treated as regular expression (true) or not (false).
	// As with global search, normally server-side processing scripts will not perform regular expression searching
	// for performance reasons on large data sets, but it is technically possible and at the discretion of your script.
	UseRegex bool
}

// DataTablesInfo represents all of the information that was requested by DataTables
type DataTablesInfo struct {
	// HasFilter Indicates there is a filter on the data to apply.  It is used to optimize generating
	// the query filters
	HasFilter bool
	// Draw counter. This is used by DataTables to ensure that the Ajax returns
	// from server-side processing requests are drawn in sequence by DataTables
	// (Ajax requests are asynchronous and thus can return out of sequence).
	// This is used as part of the draw return parameter (see below).
	Draw int
	// Start is the paging first record indicator.
	// This is the start point in the current data set (0 index based - i.e. 0 is the first record).
	Start int
	// Length is the number of records that the table can display in the current draw.
	// It is expected that the number of records returned will be equal to this number, unless the server has fewer records to return.
	//  Note that this can be -1 to indicate that all records should be returned (although that negates any benefits of server-side processing!)
	Length int
	// Searchval holds the global search value. To be applied to all columns which have searchable as true.
	Searchval string
	// UseRegex is true if the global filter should be treated as a regular expression for advanced searching.
	//  Note that normally server-side processing scripts will not perform regular expression
	//  searching for performance reasons on large data sets, but it is technically possible and at the discretion of your script.
	UseRegex bool
	// Order provides information about what columns are to be ordered in the results and which direction
	Order []OrderInfo
	// Columns provides a mapping of what fields are to be searched
	Columns []ColData
}

type DataTablesResponse struct {
	// Draw counter. This is used by DataTables to ensure that the Ajax returns
	// from server-side processing requests are drawn in sequence by DataTables
	// (Ajax requests are asynchronous and thus can return out of sequence).
	// This is used as part of the draw return parameter (see below).
	Draw int `json:"draw"`
	// Total records, before filtering (i.e. the total number of records in the database)
	RecordsTotal int `json:"recordsTotal"`
	// Total records, after filtering
	// (i.e. the total number of records after filtering has been applied - not just the number of records being returned for this page of data).
	RecordsFiltered int `json:"recordsFiltered"`
	// The data to be displayed in the table. This is an array of data source objects, one for each row, which will be used by DataTables.
	// Note that this parameter's name can be changed using the ajax option's dataSrc property.
	Data []map[string]string `json:"data"`
}

// duplicate column names
func getAllColumns(list []ColData) (allColumns []string, err error) {

	dupChecker := make(map[string]bool)

	for i, col := range list {

		// make sure col.Name not empty
		if col.Name == "" {
			err = fmt.Errorf("column %d has no name", i)
			return
		}

		// validate the column name
		if !validateColumnName(col.Name) {
			err = fmt.Errorf("column name is not valid: %s", col.Name)
			return
		}

		// check if it is a duplicate
		if dupChecker[col.Name] == true {
			err = fmt.Errorf("column %d has a duplicate name: %s", i, col.Name)
			return
		} else {
			dupChecker[col.Name] = true
		}

		// append allColumns
		allColumns = append(allColumns, col.Name)
	}

	return
}

// validate column names
func validateColumnName(name string) (isValid bool) {

	isValid, _ = regexp.MatchString("^[a-zA-Z_][a-zA-Z0-9_]*$", name)
	return

}

func mysqlEscape(sql string) string {
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]

		escape = 0

		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
			break
		case '\n': /* Must be escaped for logs */
			escape = 'n'
			break
		case '\r':
			escape = 'r'
			break
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		case '"': /* Better safe than sorry */
			escape = '"'
			break
		case '\032': /* This gives problems on Win32 */
			escape = 'Z'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}

// parseParts takes the split out parts of the field string, verifies that they are
// syntactically valid and then parses them out
//  for example columns[i][search][regex] would come in as
//       field:  'columns[2][search][regex]'
//       nameparts[0]  'columns'
//       nameparts[1]  '2]'
//       nameparts[2]  'search]'
//       nameparts[3]  'regex]'
func parseParts(field string, nameparts []string) (index int, elem1 string, elem2 string, err error) {
	defaultErr := fmt.Errorf("Invalid order[] element %v", field)
	numRegex, err := regexp.Compile("^[0-9]+]$")
	if err != nil {
		return
	}
	elemRegex, err := regexp.Compile("^[a-z]+]$")
	if err != nil {
		return
	}
	if len(nameparts) != 3 && len(nameparts) != 4 {
		err = defaultErr
		return
	}
	// Make sure it is a number followed by the closing ]
	if !numRegex.MatchString(nameparts[1]) {
		err = defaultErr
		return
	}
	// And parse it as a number to make sure
	numstr := strings.TrimSuffix(nameparts[1], "]")
	index, err = strconv.Atoi(numstr)
	if err != nil {
		return
	}
	// Check that the next index is a name token followed by a ]
	if !elemRegex.MatchString(nameparts[2]) {
		err = defaultErr
		return
	}
	// Strip off the trailing ]
	elem1 = strings.TrimSuffix(nameparts[2], "]")
	// If we had a third element, check to make sure it is also close by a ]
	if len(nameparts) == 4 {
		if !elemRegex.MatchString(nameparts[3]) {
			err = defaultErr
			return
		}
		// And trim off the ]
		elem2 = strings.TrimSuffix(nameparts[3], "]")
	}
	// Let's sanity check and make sure they aren't returning an index that is way out of range.
	// We shall assume that no more than 200 columns are being returned
	if index > 200 || index < 0 {
		err = defaultErr
	}
	return
}

// ParseDatatablesRequest checks the HTTP request to see if it corresponds
// to a datatables AJAX data request and parses the request data into
// the DataTablesInfo structure.
//
// This structure can be used by MySQLFilter and MySQLOrderby to generate a
// MySQL query to run against a database.
//
// For example assuming you are going to fill in a response structure to DataTables
// such as:
//
//   type QueryResponse struct {
//       DateAdded   time.Time
//       Status      string
//       Email       struct {
//           Name      string
//           Email     string
//       }
//   }
//   var emailQueueFields = map[string]string{
//       "DateAdded":          "t1.dateadded",
//       "Status":             "t1.status",
//       "Email.Name":         "t2.Name",
//       "Email.Email":        "t2.Email",
//   }
//
//   const baseQuery = `
//       SELECT t1.dateadded
//             ,t1.status
//             ,t2.Name
//             ,t2.Email
//       FROM infotable t1
//       LEFT JOIN usertable t2
//         ON t1.key = t2.key`
//
//       // See if we have a where clause to add to the base query
//       query := baseQuery
//       sqlPart, err := di.MySQLFilter(sqlFields)
//       // If we did have a where filter, append it.  Note that it doesn't put the " WHERE "
//       // in front because we might be doing a boolean operation.
//       if sqlPart != "" {
//           query += " WHERE " + sqlPart
//       }
//       sqlPart, err = di.MySQLOrderby(sqlFields)
//       query += sqlPart
//
// At that point you have a query that you can send straight to mySQL
//
func ParseDatatablesRequest(r *http.Request) (res *DataTablesInfo, err error) {
	var index int
	var elem string
	var elem2 string
	foundDraw := false
	res = &DataTablesInfo{}
	// Let the request parse the post values into the r.Form structure
	err = r.ParseForm()
	if err != nil {
		return
	}
	for field, value := range r.Form {
		// Remember that HTML sends us an array of values, but for datatables we only have one entry so we
		// we can shortcut and take the first element (which will be the only element) of the field.
		val0 := value[0]
		// Split out on the [ into pieces so we can see what the name is.  Note that we will have another
		// routine split out remainder of the string.
		nameparts := strings.Split(field, "[")
		switch nameparts[0] {
		case "draw":
			foundDraw = true
			res.Draw, err = strconv.Atoi(val0)
		case "start":
			res.Start, err = strconv.Atoi(val0)
		case "length":
			res.Length, err = strconv.Atoi(val0)
		case "search":
			if len(nameparts) != 2 {
				err = fmt.Errorf("Invalid search[] element %v", field)
			} else if nameparts[1] == "value]" {
				res.Searchval = val0
			} else if nameparts[1] == "regex]" {
				res.UseRegex = (val0 == "true")
			} else {
				err = fmt.Errorf("Invalid search[] element %v", field)
			}
		case "order":
			index, elem, _, err = parseParts(field, nameparts)
			if err == nil {
				// Make sure there is a spot to store this one.  Note that we may see
				// order[3][column] before we see order[0][dir]
				for len(res.Order) <= index {
					res.Order = append(res.Order, OrderInfo{})
				}
				switch elem {
				case "column":
					res.Order[index].ColNum, err = strconv.Atoi(val0)
				case "dir":
					res.Order[index].Direction = Asc
					if val0 == "desc" {
						res.Order[index].Direction = Desc
					}
				}
			}
		case "columns":
			index, elem, elem2, err = parseParts(field, nameparts)
			// First make sure we have a valid column number to work against
			if err == nil {
				// Fill up the slice to get to the spot where it is going
				// because the columns may come out of order.. I.e. we may see
				// columns[4][search][value] before we see columns[0][data]
				for len(res.Columns) <= index {
					res.Columns = append(res.Columns, ColData{})
				}
			}
			// Now fill in the field in the column slice
			switch elem {
			case "data":
				res.Columns[index].Data = val0
			case "name":
				res.Columns[index].Name = val0
			case "searchable":
				res.Columns[index].Searchable = (val0 != "false")
			case "orderable":
				res.Columns[index].Orderable = (val0 != "false")
			case "search":
				switch elem2 {
				case "value":
					res.Columns[index].Searchval = val0
				case "regex":
					res.Columns[index].UseRegex = (val0 != "false")
				}
			}
		}
		// Any errors along the way and we get out.
		if err != nil {
			return
		}
	}
	// If no Draw was specified in the request, then this isn't a datatables request and we can safely ignore it
	if !foundDraw {
		res = nil
		err = errors.New("Not a DataTables request")
	} else {
		// We have a valid datatables request.  See if we actually have any filtering
		res.HasFilter = false
		// Check the global search value to see if it has anything on it
		if res.Searchval != "" {
			// We do have a filter so note that for later
			res.HasFilter = true
			// If they ask for a regex but don't use any regular expressions, then turn off regex for efficiency
			if res.UseRegex && !strings.ContainsAny(res.Searchval, "^$.*+|[]?") {
				res.UseRegex = false
			}
			// Escape the single quotes and any escape characters and then quote the string
			res.Searchval = strings.Replace(res.Searchval, "\\", "\\\\", -1)
			res.Searchval = "'" + strings.Replace(res.Searchval, "'", "\\'", -1) + "'"
		}
		// Now we check all of the columns to see if they have search expressions
		for _, colData := range res.Columns {
			if colData.Searchval != "" {
				// We have a search expression so we remember we have a filter
				res.HasFilter = true
				// CHeck for any regular expression characters and turn off regex if not
				if colData.UseRegex && !strings.ContainsAny(colData.Searchval, "[]^$.*?+") {
					colData.UseRegex = false
				}
				// Escape the single quotes and any escape characters and then quote the string
				colData.Searchval = strings.Replace(colData.Searchval, "\\", "\\\\", -1)
				colData.Searchval = "'" + strings.Replace(colData.Searchval, "'", "\\'", -1) + "'"
			}
		}
	}
	return
}

func (di *DataTablesInfo) SetDbX(db *sqlx.DB) {
	Dbx = *db
}

func (di *DataTablesInfo) fetchDataForResponse(tableName string) (response DataTablesResponse, err error) {

	// validate column data
	allColumns, err := getAllColumns(di.Columns)
	if err != nil {
		return
	}

	// generate WHERE condition
	whereClause, err := di.generateWhereClause()
	if err != nil {
		return
	}

	// generate ORDER BY condition
	orderClause, err := di.generateOrderClause()
	if err != nil {
		return
	}

	// query to find total rows
	queryTotal := fmt.Sprintf(
		"SELECT COUNT(present) AS count FROM %s WHERE present=1",
		tableName,
	)
	// view_portgroup doesnt have "present" column
	if tableName == "view_portgroup" {
		queryTotal = fmt.Sprintf(
			"SELECT COUNT(*) AS count FROM %s",
			tableName,
		)
	}

	// query to find return all filtered data
	queryFiltered := fmt.Sprintf(
		"SELECT SQL_CALC_FOUND_ROWS %s FROM %s %s %s LIMIT %d, %d",
		strings.Join(allColumns, ","),
		tableName,
		whereClause,
		orderClause,
		di.Start,
		di.Length,
	)

	//log.Debugf("sql query: %s", queryFiltered)

	// query to find total rows filtered
	queryFilteredCount := "SELECT FOUND_ROWS() AS rows;"

	// populate the resonse
	err = querySingleRow(queryTotal, &response.RecordsTotal)
	if err != nil {
		return
	}

	err = populateData(queryFiltered, &response)
	if err != nil {
		return
	}

	err = querySingleRow(queryFilteredCount, &response.RecordsFiltered)
	if err != nil {
		return
	}

	response.Draw = di.Draw

	return

}

// this will generate the WHERE condition of query
func (di *DataTablesInfo) generateWhereClause() (whereClause string, err error) {

	extra := "WHERE ("
	for _, colData := range di.Columns {

		if colData.Searchable {
			// If we have a global search val, generate a match against the global value for this field
			if di.Searchval != "" {
				// For wildcards we have to generate a REGEXP request
				if di.UseRegex {
					whereClause += extra + colData.Name + " REGEX " + di.Searchval
					extra = " OR "
				} else {
					// In the special case where we have a top level non wild card search value we want
					// to gang all the fields together into a single match string
					whereClause += extra + "MATCH(" + colData.Name + ") AGAINST(" + di.Searchval + ")"
					extra = " OR "
				}
			}
			// See if we have a search value specific for this individual element
			if colData.Searchval != "" {
				if colData.UseRegex {
					whereClause += extra + colData.Name + " REGEXP '" + colData.Searchval + "'"
				} else {
					whereClause += extra + colData.Name + " LIKE '%" + mysqlEscape(colData.Searchval) + "%'"
				}
				extra = " AND "
			}
		}
	}

	if whereClause != "" {
		whereClause += ")"
	}

	return
}

// this will generate the ORDER BY part of query
func (di *DataTablesInfo) generateOrderClause() (orderClause string, err error) {

	extra := "ORDER BY "
	// Go through the list of requested items to order
	for _, orderItem := range di.Order {
		// Make sure that the column is in range
		if orderItem.ColNum >= len(di.Columns) {
			err = fmt.Errorf("Datatables Request order column %v out of range %v of columns", orderItem.ColNum, len(di.Columns))
			return
		}
		// Get the data for that column and figure out if the name is one of the fields that we
		// allow in the table
		colData := di.Columns[orderItem.ColNum]
		// Make sure we can actually order on the column (in theory this will never happen)
		if !colData.Orderable {
			err = fmt.Errorf("Datatables requested ordering on non-orderable column %v", colData.Data)
			return
		}
		// We have the column in the database, add it to the order by query that we are generating
		// The first time we have " ORDER BY " in the extra string, subsequent times we get a simple ","
		// which allows us to build up the string without backtracking to remove characters
		orderClause += extra + colData.Name
		if orderItem.Direction == Desc {
			orderClause += " DESC"
		}
		extra = ","
	}
	// If for some reason we got to the end with no columns, then we give them the order by the first item
	if orderClause == "" {
		orderClause = extra + "1"
	}

	return
}

// generic single row query
func querySingleRow(query string, arg interface{}) (err error) {

	err = Dbx.QueryRow(query).Scan(arg)
	return
}

// this will populate the data part of the response
func populateData(query string, response *DataTablesResponse) (err error) {

	rows, err := Dbx.Query(query)
	if err != nil {
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		return
	}

	// create dynamic []map[string]string
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		temp := map[string]string{}

		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			temp[col] = fmt.Sprintf("%s", v)
		}

		// add this to the response
		response.Data = append(response.Data, temp)
	}

	return
}
