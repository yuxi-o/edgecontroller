// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React from "react";
import TableRow from "@material-ui/core/TableRow";
import TableCell from "@material-ui/core/TableCell";
import Button from "@material-ui/core/Button";
import { Link } from "react-router-dom";

class EnhancedTableRow extends React.Component {

  render() {
    const { tableData } = this.props;
    const isSelected = this.props.isSelected;

    const editLink = (props) => <Link to={tableData.editUrl} {...props} />;

    return (
      <TableRow
        hover
        aria-checked={isSelected}
        tabIndex={-1}
        key={tableData.id}
        selected={isSelected}
      >
        {Object.keys(tableData).map(key => {
          return key !== 'editUrl' ? (
            <TableCell align="left">{tableData[key]}</TableCell>
          ) : "";
        })}
        <TableCell align="left">
          <Button
            style={{ marginRight: '15px' }}
            variant="contained"
            color="primary"
            component={editLink}
            className={'classes.button'}
          >
            Edit
          </Button>
        </TableCell>
      </TableRow>
    );
  }
};

export default EnhancedTableRow
