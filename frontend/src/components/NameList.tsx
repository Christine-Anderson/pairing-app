import React, { useState } from 'react';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import Button from '@mui/material/Button';
import Paper from '@mui/material/Paper';
import FormControl from '@mui/material/FormControl';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import Typography from '@mui/material/Typography';

interface NameListProps {
    names: string[];
}

const NameList = ({names}: NameListProps) => {
    const [selectedValues, setSelectedValues] = useState<{ [key: string]: string[] }>({});

    const handleChange = (name: string, event: React.ChangeEvent<{ value: unknown }>) => {
        setSelectedValues(prev => ({
            ...prev,
            [name]: event.target.value as string[]
        }));
    };

    const handleButtonClick = () => {
        // todo form submission
        console.log("Button clicked!");
    };

    return (
        <Paper elevation={3} style={{ padding: '1rem', maxWidth: '20rem', margin: 'auto' }}>
            <Typography variant="subtitle1" align="center" gutterBottom>
                Please include any restrictions in who the participants can be paired with.
            </Typography>
            <List>
                {names.map((name, index) => (
                    <ListItem key={index}>
                        <ListItemText primary={name} />
                        <FormControl style={{ minWidth: '12rem', marginLeft: '1rem' }}>
                            <Select
                                multiple
                                value={selectedValues[name] || []}
                                onChange={(event) => handleChange(name, event as React.ChangeEvent<{ value: unknown }>)}
                                renderValue={(selected) => (selected).join(', ')}
                                style={{ minWidth: '9em' }}
                            >
                                {names.filter(option => option !== name).map((option, idx) => (
                                    <MenuItem key={idx} value={option}>{option}</MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                    </ListItem>
                ))}
            </List>
            <div style={{ textAlign: 'center', marginTop: '1em' }}>
                <Button variant="contained" color="primary" onClick={handleButtonClick}>
                    Pair
                </Button>
            </div>
        </Paper>
    );
};

export default NameList;
