import React, { useState } from "react";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import Button from "@mui/material/Button";
import Paper from "@mui/material/Paper";
import FormControl from "@mui/material/FormControl";
import Select from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import Typography from "@mui/material/Typography";
import LinearProgress from '@mui/material/LinearProgress';
import { GroupMember } from "../types";

interface NameListProps {
    groupMembers: GroupMember[];
    isLoading: boolean;
}

const NameList = ({groupMembers, isLoading}: NameListProps) => {
    const [selectedValues, setSelectedValues] = useState<{ [key: string]: string[] }>({});
    const handleChange = (memberId: string, event: React.ChangeEvent<{ value: unknown }>) => {
        const selectedIds = event.target.value as string[];
        setSelectedValues(prev => ({
            ...prev,
            [memberId]: selectedIds,
        }));
    };
    const handleButtonClick = () => {
        // todo form submission
        console.log("Selected IDs:", selectedValues);
    };
    return (
        <Paper elevation={3} style={{ padding: "1rem", maxWidth: "20rem", margin: "auto" }}>
            { isLoading && <LinearProgress/> }
            {
                groupMembers && 
                <Typography variant="subtitle1" align="center" gutterBottom>
                    Please include any restrictions in who the participants can be paired with.
                </Typography>
            }
            <List>
                {groupMembers && groupMembers.map(({memberId, name}, index) => (
                    <ListItem key={index}>
                        <ListItemText primary={name} />
                        <FormControl style={{ minWidth: "12rem", marginLeft: "1rem" }}>
                            <Select
                                multiple
                                value={selectedValues[memberId] || []}
                                onChange={(event) => handleChange(memberId, event as React.ChangeEvent<HTMLInputElement>)}
                                renderValue={(selected) => {
                                    const selectedNames = (selected).map(memberId => {
                                        const selectedPerson = groupMembers.find(item => item.memberId === memberId);
                                        return selectedPerson ? selectedPerson.name : "";
                                    });
                                    return selectedNames.join(", ");
                                }}
                                style={{ minWidth: "9em" }}
                            >
                                {groupMembers
                                    .filter(member => member.memberId !== memberId)
                                    .map(({ memberId, name }) => (
                                        <MenuItem key={memberId} value={memberId}>
                                            {name}
                                        </MenuItem>
                                    ))}
                            </Select>
                        </FormControl>
                    </ListItem>
                ))}
            </List>
            <div style={{ textAlign: "center", marginTop: "1em" }}>
                <Button variant="contained" color="primary" onClick={handleButtonClick}>
                    Submit
                </Button>
            </div>
        </Paper>
    );
};
export default NameList;