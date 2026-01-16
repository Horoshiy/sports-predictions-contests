import React, { useMemo } from 'react'
import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Box, Grid, FormControl, InputLabel, Select, MenuItem } from '@mui/material'
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { matchSchema, type MatchFormData } from '../../utils/sports-validation'
import { useLeagues, useTeams } from '../../hooks/use-sports'
import type { Match } from '../../types/sports.types'

interface MatchFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: MatchFormData) => void
  match?: Match | null
  loading?: boolean
}

export const MatchForm: React.FC<MatchFormProps> = ({ open, onClose, onSubmit, match, loading = false }) => {
  const isEditing = !!match
  const { data: leaguesData } = useLeagues({ pagination: { page: 1, limit: 100 }, activeOnly: true })
  const { data: teamsData } = useTeams({ pagination: { page: 1, limit: 200 }, activeOnly: true })

  const defaultValues = React.useMemo(() => ({
    leagueId: match?.leagueId || 0,
    homeTeamId: match?.homeTeamId || 0,
    awayTeamId: match?.awayTeamId || 0,
    scheduledAt: match?.scheduledAt ? new Date(match.scheduledAt) : new Date(),
    status: match?.status || 'scheduled',
    homeScore: match?.homeScore || 0,
    awayScore: match?.awayScore || 0,
    resultData: match?.resultData || '',
  }), [match])

  const { control, handleSubmit, reset, watch, formState: { errors, isValid } } = useForm<MatchFormData>({
    resolver: zodResolver(matchSchema),
    defaultValues,
    mode: 'onChange',
  })

  const selectedLeagueId = watch('leagueId')
  const selectedLeague = leaguesData?.leagues?.find(l => l.id === selectedLeagueId)

  const filteredTeams = useMemo(() => {
    if (!selectedLeague || !teamsData?.teams) return []
    return teamsData.teams.filter(t => t.sportId === selectedLeague.sportId)
  }, [selectedLeague, teamsData])

  // Reset team selections when league changes (different sport may have different teams)
  const prevLeagueId = React.useRef(selectedLeagueId)
  React.useEffect(() => {
    if (!isEditing && prevLeagueId.current !== selectedLeagueId && prevLeagueId.current !== 0) {
      setValue('homeTeamId', 0)
      setValue('awayTeamId', 0)
    }
    prevLeagueId.current = selectedLeagueId
  }, [selectedLeagueId, isEditing, setValue])

  React.useEffect(() => {
    reset(defaultValues)
  }, [defaultValues, reset])

  const handleClose = () => {
    reset()
    onClose()
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>{isEditing ? 'Edit Match' : 'Create Match'}</DialogTitle>
      <form onSubmit={handleSubmit(onSubmit)}>
        <DialogContent>
          <LocalizationProvider dateAdapter={AdapterDateFns}>
            <Box sx={{ mt: 1 }}>
              <Grid container spacing={2}>
                <Grid item xs={12}>
                  <Controller
                    name="leagueId"
                    control={control}
                    render={({ field }) => (
                      <FormControl fullWidth required error={!!errors.leagueId}>
                        <InputLabel>League</InputLabel>
                        <Select {...field} label="League" disabled={loading}>
                          {leaguesData?.leagues?.map(l => <MenuItem key={l.id} value={l.id}>{l.name}</MenuItem>)}
                        </Select>
                      </FormControl>
                    )}
                  />
                </Grid>
                <Grid item xs={6}>
                  <Controller
                    name="homeTeamId"
                    control={control}
                    render={({ field }) => (
                      <FormControl fullWidth required error={!!errors.homeTeamId}>
                        <InputLabel>Home Team</InputLabel>
                        <Select {...field} label="Home Team" disabled={loading || !selectedLeagueId}>
                          {filteredTeams.map(t => <MenuItem key={t.id} value={t.id}>{t.name}</MenuItem>)}
                        </Select>
                      </FormControl>
                    )}
                  />
                </Grid>
                <Grid item xs={6}>
                  <Controller
                    name="awayTeamId"
                    control={control}
                    render={({ field }) => (
                      <FormControl fullWidth required error={!!errors.awayTeamId}>
                        <InputLabel>Away Team</InputLabel>
                        <Select {...field} label="Away Team" disabled={loading || !selectedLeagueId}>
                          {filteredTeams.map(t => <MenuItem key={t.id} value={t.id}>{t.name}</MenuItem>)}
                        </Select>
                        {errors.awayTeamId && <Box sx={{ color: 'error.main', fontSize: '0.75rem', mt: 0.5 }}>{errors.awayTeamId.message}</Box>}
                      </FormControl>
                    )}
                  />
                </Grid>
                <Grid item xs={12}>
                  <Controller
                    name="scheduledAt"
                    control={control}
                    render={({ field }) => (
                      <DateTimePicker {...field} label="Scheduled At" disabled={loading} slotProps={{ textField: { fullWidth: true, error: !!errors.scheduledAt, helperText: errors.scheduledAt?.message } }} />
                    )}
                  />
                </Grid>
                {isEditing && (
                  <>
                    <Grid item xs={12}>
                      <Controller
                        name="status"
                        control={control}
                        render={({ field }) => (
                          <FormControl fullWidth>
                            <InputLabel>Status</InputLabel>
                            <Select {...field} label="Status" disabled={loading}>
                              <MenuItem value="scheduled">Scheduled</MenuItem>
                              <MenuItem value="live">Live</MenuItem>
                              <MenuItem value="finished">Finished</MenuItem>
                              <MenuItem value="cancelled">Cancelled</MenuItem>
                              <MenuItem value="postponed">Postponed</MenuItem>
                            </Select>
                          </FormControl>
                        )}
                      />
                    </Grid>
                    <Grid item xs={6}>
                      <Controller name="homeScore" control={control} render={({ field }) => (
                        <TextField {...field} label="Home Score" type="number" fullWidth disabled={loading} onChange={(e) => field.onChange(parseInt(e.target.value) || 0)} />
                      )} />
                    </Grid>
                    <Grid item xs={6}>
                      <Controller name="awayScore" control={control} render={({ field }) => (
                        <TextField {...field} label="Away Score" type="number" fullWidth disabled={loading} onChange={(e) => field.onChange(parseInt(e.target.value) || 0)} />
                      )} />
                    </Grid>
                  </>
                )}
              </Grid>
            </Box>
          </LocalizationProvider>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} disabled={loading}>Cancel</Button>
          <Button type="submit" variant="contained" disabled={loading || !isValid}>
            {loading ? 'Saving...' : isEditing ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}

export default MatchForm
