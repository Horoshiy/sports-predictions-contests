import React from 'react'
import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Box, Grid, FormControl, InputLabel, Select, MenuItem, FormHelperText } from '@mui/material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { leagueSchema, type LeagueFormData, generateSlug } from '../../utils/sports-validation'
import { useSports } from '../../hooks/use-sports'
import type { League } from '../../types/sports.types'

interface LeagueFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: LeagueFormData) => void
  league?: League | null
  loading?: boolean
}

export const LeagueForm: React.FC<LeagueFormProps> = ({ open, onClose, onSubmit, league, loading = false }) => {
  const isEditing = !!league
  const { data: sportsData } = useSports({ pagination: { page: 1, limit: 100 }, activeOnly: true })

  const defaultValues = React.useMemo(() => ({
    sportId: league?.sportId || 0,
    name: league?.name || '',
    slug: league?.slug || '',
    country: league?.country || '',
    season: league?.season || '',
  }), [league])

  const { control, handleSubmit, reset, watch, setValue, formState: { errors, isValid } } = useForm<LeagueFormData>({
    resolver: zodResolver(leagueSchema),
    defaultValues,
    mode: 'onChange',
  })

  const name = watch('name')
  const slug = watch('slug')

  React.useEffect(() => {
    if (!isEditing && name && !slug) {
      setValue('slug', generateSlug(name))
    }
  }, [name, slug, isEditing, setValue])

  React.useEffect(() => {
    reset(defaultValues)
  }, [defaultValues, reset])

  const handleClose = () => {
    reset()
    onClose()
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>{isEditing ? 'Edit League' : 'Create League'}</DialogTitle>
      <form onSubmit={handleSubmit(onSubmit)}>
        <DialogContent>
          <Box sx={{ mt: 1 }}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Controller
                  name="sportId"
                  control={control}
                  render={({ field }) => (
                    <FormControl fullWidth required error={!!errors.sportId}>
                      <InputLabel>Sport</InputLabel>
                      <Select {...field} label="Sport" disabled={loading}>
                        {sportsData?.sports?.map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}
                      </Select>
                      {errors.sportId && <FormHelperText>{errors.sportId.message}</FormHelperText>}
                    </FormControl>
                  )}
                />
              </Grid>
              <Grid item xs={12}>
                <Controller name="name" control={control} render={({ field }) => (
                  <TextField {...field} label="Name" fullWidth required error={!!errors.name} helperText={errors.name?.message} disabled={loading} />
                )} />
              </Grid>
              <Grid item xs={12}>
                <Controller name="slug" control={control} render={({ field }) => (
                  <TextField {...field} label="Slug" fullWidth error={!!errors.slug} helperText={errors.slug?.message || 'Auto-generated'} disabled={loading} />
                )} />
              </Grid>
              <Grid item xs={6}>
                <Controller name="country" control={control} render={({ field }) => (
                  <TextField {...field} label="Country" fullWidth error={!!errors.country} helperText={errors.country?.message} disabled={loading} />
                )} />
              </Grid>
              <Grid item xs={6}>
                <Controller name="season" control={control} render={({ field }) => (
                  <TextField {...field} label="Season" fullWidth error={!!errors.season} helperText={errors.season?.message} disabled={loading} />
                )} />
              </Grid>
            </Grid>
          </Box>
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

export default LeagueForm
