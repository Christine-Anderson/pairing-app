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

interface GroupMember {
    id: string;
    name: string;
}

interface NameListProps {
    groupMembers: GroupMember[];
}

const NameList = ({groupMembers}: NameListProps) => {
    const [selectedValues, setSelectedValues] = useState<{ [key: string]: string[] }>({});
    const handleChange = (id: string, event: React.ChangeEvent<{ value: unknown }>) => {
        const selectedIds = event.target.value as string[];
        setSelectedValues(prev => ({
            ...prev,
            [id]: selectedIds,
        }));
    };
    const handleButtonClick = () => {
        // todo form submission
        console.log("Selected IDs:", selectedValues);
    };
    return (
        <Paper elevation={3} style={{ padding: "1rem", maxWidth: "20rem", margin: "auto" }}>
            <Typography variant="subtitle1" align="center" gutterBottom>
                Please include any restrictions in who the participants can be paired with.
            </Typography>
            <List>
                {groupMembers && groupMembers.map(({id, name}, index) => (
                    <ListItem key={index}>
                        <ListItemText primary={name} />
                        <FormControl style={{ minWidth: "12rem", marginLeft: "1rem" }}>
                            <Select
                                multiple
                                value={selectedValues[id] || []}
                                onChange={(event) => handleChange(id, event as React.ChangeEvent<HTMLInputElement>)}
                                renderValue={(selected) => {
                                    const selectedNames = (selected).map(id => {
                                        const selectedPerson = groupMembers.find(item => item.id === id);
                                        return selectedPerson ? selectedPerson.name : "";
                                    });
                                    return selectedNames.join(", ");
                                }}
                                style={{ minWidth: "9em" }}
                            >
                                {groupMembers
                                    .filter(member => member.id !== id)
                                    .map(({ id, name }) => (
                                        <MenuItem key={id} value={id}>
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