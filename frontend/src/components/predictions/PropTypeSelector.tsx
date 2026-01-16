import React from 'react'
import {
  Box,
  Card,
  CardContent,
  Typography,
  FormControl,
  RadioGroup,
  FormControlLabel,
  Radio,
  TextField,
  Chip,
  IconButton,
  Stack,
} from '@mui/material'
import { Add as AddIcon, Delete as DeleteIcon } from '@mui/icons-material'
import type { PropType } from '../../types/props.types'
import type { PropPredictionFormData } from '../../utils/prediction-validation'

interface PropTypeSelectorProps {
  propTypes: PropType[]
  selectedProps: PropPredictionFormData[]
  onPropsChange: (props: PropPredictionFormData[]) => void
  homeTeam: string
  awayTeam: string
  disabled?: boolean
}

export const PropTypeSelector: React.FC<PropTypeSelectorProps> = ({
  propTypes,
  selectedProps,
  onPropsChange,
  homeTeam,
  awayTeam,
  disabled = false,
}) => {
  const addProp = (propType: PropType) => {
    const newProp: PropPredictionFormData = {
      propTypeId: propType.id,
      propSlug: propType.slug,
      line: propType.defaultLine ?? undefined,
      selection: '',
      pointsValue: propType.pointsCorrect,
    }
    onPropsChange([...selectedProps, newProp])
  }

  const removeProp = (index: number) => {
    onPropsChange(selectedProps.filter((_, i) => i !== index))
  }

  const updateProp = (index: number, updates: Partial<PropPredictionFormData>) => {
    const updated = [...selectedProps]
    updated[index] = { ...updated[index], ...updates }
    onPropsChange(updated)
  }

  const getSelectionOptions = (propType: PropType) => {
    switch (propType.valueType) {
      case 'over_under':
        return [
          { value: 'over', label: 'Over' },
          { value: 'under', label: 'Under' },
        ]
      case 'yes_no':
        return [
          { value: 'yes', label: 'Yes' },
          { value: 'no', label: 'No' },
        ]
      case 'team_select':
        return [
          { value: 'home', label: homeTeam },
          { value: 'away', label: awayTeam },
        ]
      default:
        return []
    }
  }

  const availablePropTypes = propTypes.filter(
    pt => !selectedProps.some(sp => sp.propTypeId === pt.id)
  )

  const groupedPropTypes = React.useMemo(() => 
    availablePropTypes.reduce((acc, pt) => {
      if (!acc[pt.category]) acc[pt.category] = []
      acc[pt.category].push(pt)
      return acc
    }, {} as Record<string, PropType[]>),
    [availablePropTypes]
  )

  return (
    <Box>
      {selectedProps.length > 0 && (
        <Box sx={{ mb: 3 }}>
          <Typography variant="subtitle2" gutterBottom>
            Your Props ({selectedProps.length})
          </Typography>
          <Stack spacing={2}>
            {selectedProps.map((prop, index) => {
              const propType = propTypes.find(pt => pt.id === prop.propTypeId)
              if (!propType) return null

              return (
                <Card key={prop.propTypeId} variant="outlined">
                  <CardContent sx={{ py: 1.5, '&:last-child': { pb: 1.5 } }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                      <Box sx={{ flex: 1 }}>
                        <Typography variant="body2" fontWeight="medium">
                          {propType.name}
                          <Chip label={`+${propType.pointsCorrect} pts`} size="small" color="primary" sx={{ ml: 1 }} />
                        </Typography>
                        
                        {propType.valueType === 'over_under' && (
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mt: 1 }}>
                            <TextField
                              label="Line"
                              type="number"
                              size="small"
                              value={prop.line ?? ''}
                              onChange={(e) => updateProp(index, { 
                                line: e.target.value ? parseFloat(e.target.value) : undefined 
                              })}
                              disabled={disabled}
                              sx={{ width: 100 }}
                              inputProps={{ step: 0.5 }}
                            />
                            <FormControl size="small">
                              <RadioGroup
                                row
                                value={prop.selection}
                                onChange={(e) => updateProp(index, { selection: e.target.value })}
                              >
                                {getSelectionOptions(propType).map(opt => (
                                  <FormControlLabel
                                    key={opt.value}
                                    value={opt.value}
                                    control={<Radio size="small" />}
                                    label={opt.label}
                                    disabled={disabled}
                                  />
                                ))}
                              </RadioGroup>
                            </FormControl>
                          </Box>
                        )}

                        {(propType.valueType === 'yes_no' || propType.valueType === 'team_select') && (
                          <FormControl size="small" sx={{ mt: 1 }}>
                            <RadioGroup
                              row
                              value={prop.selection}
                              onChange={(e) => updateProp(index, { selection: e.target.value })}
                            >
                              {getSelectionOptions(propType).map(opt => (
                                <FormControlLabel
                                  key={opt.value}
                                  value={opt.value}
                                  control={<Radio size="small" />}
                                  label={opt.label}
                                  disabled={disabled}
                                />
                              ))}
                            </RadioGroup>
                          </FormControl>
                        )}
                      </Box>
                      <IconButton size="small" onClick={() => removeProp(index)} disabled={disabled}>
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    </Box>
                  </CardContent>
                </Card>
              )
            })}
          </Stack>
        </Box>
      )}

      <Typography variant="subtitle2" gutterBottom>
        Add Props
      </Typography>
      {Object.entries(groupedPropTypes).map(([category, types]) => (
        <Box key={category} sx={{ mb: 2 }}>
          <Typography variant="caption" color="text.secondary" sx={{ textTransform: 'capitalize' }}>
            {category} Props
          </Typography>
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mt: 0.5 }}>
            {types.map(propType => (
              <Chip
                key={propType.id}
                label={propType.name}
                onClick={() => addProp(propType)}
                onDelete={() => addProp(propType)}
                deleteIcon={<AddIcon />}
                variant="outlined"
                disabled={disabled}
                size="small"
              />
            ))}
          </Box>
        </Box>
      ))}

      {availablePropTypes.length === 0 && selectedProps.length > 0 && (
        <Typography variant="body2" color="text.secondary">
          All available props selected
        </Typography>
      )}
    </Box>
  )
}

export default PropTypeSelector
