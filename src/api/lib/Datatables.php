<?php
namespace Ozdemir\Datatables;

use Ozdemir\Datatables\DB\DatabaseInterface;

class Datatables {

    private $db;
    private $data;
    private $recordstotal;
    private $recordsfiltered;
    private $columns;
    private $edit;

    function __construct(DatabaseInterface $db)
    {
        $this->db = $db->connect();
        $this->input = isset($_POST["draw"]) ? $_POST : $_GET;
    }

    public function query($query)
    {
<<<<<<< HEAD
        $this->setColumns($query);
        $columns = implode(", ", $this->columns);
        $query = str_replace(";", "", $query);
        $sql = "Select $columns from ($query)t";
        $this->recordstotal = $this->getCount($sql); // unfiltered data count is here.

        // filtering via global search
        $search = "";
        $globalsearch = $this->input('search')['value'];

        if ($globalsearch <> "")
        {
            $search = " WHERE (";
            foreach ($this->columns as $column)
            {
                $lookfor[] = $column . " LIKE " . $this->db->escape('%' . $globalsearch . '%') . "";
            }
            $search .= implode(" OR ", $lookfor) . ")";
        }

        // todo: Individual column filtering

        $this->recordsfiltered = $this->getCount($sql . $search);  // filtered data count is here.

        $this->data = $this->db->query($sql . $search . $this->orderby() . $this->limit());

    }

    function setColumns($query)
    {
        preg_match_all("/SELECT([\s\S]*?)FROM([\s\S]*?)/i", $query, $columns);
        $columns = $this->explode(",", $columns[1][0]);
        $columns = preg_replace("/(.*)\s+as\s+(.*)/i", "$2", $columns);
        $columns = preg_replace("/(.+)(\([^()]+\))?\s+(.+)/i", "$3", $columns);
        $columns = preg_replace('/[\s"\'`]+/', '', $columns);
        $this->columns = preg_replace("/([\w\-]*)\.([\w\-]*)/", "$2", $columns);
    }

    private function getCount($query)
    {
        return $this->db->count($query);
=======
        $this->columns = $this->setcolumns($query);
        $columns = implode(", ", $this->columns);
        $query = rtrim($query, "; ");
        $sql = "Select $columns from ($query)t";

        $this->recordstotal = $this->db->count($sql); // unfiltered data count is here.
        $where = $this->filter();
        $this->recordsfiltered = $this->db->count($sql . $where);  // filtered data count is here.

        $this->data = $this->db->query($sql . $where . $this->orderby() . $this->limit());

        return $this;
    }

    private function filter()
    {
        $search = "";

        $filterglobal = $this->filterglobal();
        $filterindividual = $this->filterindividual();

        if ($filterglobal <> null)
        {
            $search .= $filterglobal;
        }

        if ($filterindividual <> null)
        {
            if ($filterglobal <> null)
            {
                $search .= " AND ";
            }
            $search .= $filterindividual;
        }

        if ($filterindividual || $filterglobal)
        {
            $search = " WHERE " . $search;

            return $search;
        }

        return null;
    }

    private function filterglobal()
    {
        $searchinput = $this->input('search')['value'];
        $allcolumns = $this->input('columns');

        if ($searchinput <> "")
        {
            $search = " (";

            foreach ($this->columns as $key => $column)
            {
                if ($allcolumns[ $key ]['searchable'] == 'true')
                {
                    $lookfor[] = $column . " LIKE " . $this->db->escape('%' . $searchinput . '%') . "";
                }
            }

            $search .= implode(" OR ", $lookfor) . ")";

            return $search;
        }

        return null;
    }

    private function filterindividual()
    {
        $allcolumns = $this->input('columns');

        $search = " (";
        $lookfor = [];

        foreach ($allcolumns as $key)
        {
            if ($key['search']['value'] <> "" and $key['searchable'] == 'true')
            {
                $lookfor[] = $this->column($key['data']) . " LIKE " . $this->db->escape('%' . $key['search']['value'] . '%') . "";
            }
        }

        if (count($lookfor) > 0)
        {
            $search .= implode(" AND ", $lookfor) . ")";

            return $search;
        }

        return null;
    }

    private function setcolumns($query)
    {
        preg_match_all("/SELECT([\s\S]*?)((\s*)FROM(?![\s\S]*\)))([\s\S]*?)/i", $query, $columns);

        $columns = $this->explode(",", $columns[1][0]);

        // gets alias of the table -> 'table.column as col' or 'table.column col' to 'col'
        $regex[0] = "/(.*)\s+as\s+(.*)/i";
        $regex[1] = "/.+(\([^()]+\))?\s+(.+)/i";
        // wipe unwanted characters => '`" and space
        $regex[2] = '/[\s"\'`]+/';
        // if there is no alias, return column name -> table.column to column
        $regex[3] = "/([\w\-]*)\.([\w\-]*)/";

        foreach ($regex as $pattern)
        {
            $columns = preg_replace($pattern, "$2", $columns);
        }

        return $columns;
>>>>>>> origin/master
    }

    private function limit()
    {
<<<<<<< HEAD
        $skip = (integer) $this->input('start');
        $take = (integer) $this->input('length');
=======
        $take = 10;
        $skip = (integer) $this->input('start');

        if ($this->input('length'))
        {
            $take = (integer) $this->input('length');
        }
>>>>>>> origin/master

        if ($take == - 1)
        {
            return null;
        }

        return " LIMIT $skip, $take ";
    }

    private function orderby()
    {
        $dtorders = $this->input('order');
        $orders = " ORDER BY ";
        $dir = ['asc' => 'asc', 'desc' => 'desc'];

<<<<<<< HEAD
        if (is_array($dtorders))
        {
            foreach ($dtorders as $order)
            {
                $takeorders[] = $this->columns[ $order['column'] ] . " " . $dir[ $order['dir'] ];
            }
            $orders .= implode(",", $takeorders);
        } else
        {   // nothing to order use default
            $orders .= $this->columns[0] . " asc";
        }

        return $orders;
=======
        if ( ! is_array($dtorders))
        {
            return $orders . $this->columns[0] . " asc";  // default
        }

        foreach ($dtorders as $order)
        {
            $takeorders[] = $this->columns[ $order['column'] ] . " " . $dir[ $order['dir'] ];
        }

        return $orders . implode(",", $takeorders);
>>>>>>> origin/master
    }

    public function generate()
    {
        $formatted_data = [];

        foreach ($this->data as $key => $row)
        {
            // editing columns..
            if (count($this->edit) > 0)
            {
<<<<<<< HEAD
                foreach ($this->edit as $edit_job => $edit_column)
                {
                    foreach ($edit_column as $closure)
                    {
                        $row[ $edit_job ] = $this->exec_replace($closure, $row);
=======
                foreach ($this->edit as $edit_column => $edit_job)
                {
                    foreach ($edit_job as $closure)
                    {
                        $row[ $edit_column ] = $this->exec_replace($closure, $row);
>>>>>>> origin/master
                    }
                }
            }

<<<<<<< HEAD
            // Check datatables if it uses column names as data keys or not.
            $formatted_data[] = $this->isIndexed($row);

=======
            $formatted_data[] = $this->isIndexed($row);
>>>>>>> origin/master
        }

        return $this->response($formatted_data);
    }

    public function edit($column, $closure)
    {
        $this->edit[ $column ][] = $closure;

        return $this;
    }

    private function input($input)
    {
        if (isset($this->input[ $input ]))
        {
            return $this->input[ $input ];
        }

        return false;
    }

<<<<<<< HEAD
=======
    private function column($input)
    {
        if (is_numeric($input))
        {
            return $this->columns[ $input ];
        }

        return $input;
    }

>>>>>>> origin/master
    private function exec_replace($closure, $row_data)
    {
        // if this is a closure function, return calculated data.
        if (is_object($closure))
        {
            if (get_class($closure) == 'Closure')
            {
                return $closure($row_data);
            }
        }

        return false;
    }

    private function response($data)
    {
        $response['draw'] = $this->input('draw');
        $response['recordsTotal'] = $this->recordstotal;
        $response['recordsFiltered'] = $this->recordsfiltered;
        $response['data'] = $data;

        return json_encode($response);
    }

<<<<<<< HEAD
    private function isIndexed($row) // if data source uses associative keys or index
=======
    private function isIndexed($row) // if data source uses associative keys or index number
>>>>>>> origin/master
    {
        $column = $this->input('columns');
        if (is_numeric($column[0]['data']))
        {
            return array_values($row);
        }

        return $row;
    }

    private function balanceChars($str, $open, $close)
    {
        $openCount = substr_count($str, $open);
        $closeCount = substr_count($str, $close);
        $retval = $openCount - $closeCount;

        return $retval;
    }

    private function explode($delimiter, $str, $open = '(', $close = ')')
    {
        $retval = array();
        $hold = array();
        $balance = 0;
        $parts = explode($delimiter, $str);
        foreach ($parts as $part)
        {
            $hold[] = $part;
            $balance += $this->balanceChars($part, $open, $close);
            if ($balance < 1)
            {
                $retval[] = implode($delimiter, $hold);
                $hold = array();
                $balance = 0;
            }
        }
        if (count($hold) > 0)
        {
            $retval[] = implode($delimiter, $hold);
        }

        return $retval;
    }
}