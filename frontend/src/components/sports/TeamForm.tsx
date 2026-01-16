import React from 'react'
import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Box, Grid, FormControl, InputLabel, Select, MenuItem, Avatar, FormHelperText } from '@mui/material'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { teamSchema, type TeamFormData, generateSlug } from '../../utils/sports-validation'
import { useSports } from '../../hooks/use-sports'
import type { Team } from '../../types/sports.types'

interface TeamFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: TeamFormData) => void
  team?: Team | null
  loading?: boolean
}

export const TeamForm: React.FC<TeamFormProps> = ({ open, onClose, onSubmit, team, loading = false }) => {
  const isEditing = !!team
  const { data: sportsData } = useSports({ pagination: { page: 1, limit: 100 }, activeOnly: true })

  const defaultValues = React.useMemo(() => ({
    sportId: team?.sportId || 0,
    name: team?.name || '',
    slug: team?.slug || '',
    shortName: team?.shortName || '',
    logoUrl: team?.logoUrl || '',
    country: team?.country || '',
  }), [team])

  const { control, handleSubmit, reset, watch, setValue, formState: { errors, isValid } } = useForm<TeamFormData>({
    resolver: zodResolver(teamSchema),
    defaultValues,
    mode: 'onChange',
  })

  const name = watch('name')
  const slug = watch('slug')
  const logoUrl = watch('logoUrl')

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
      <DialogTitle>{isEditing ? 'Edit Team' : 'Create Team'}</DialogTitle>
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
              <Grid item xs={8}>
                <Controller name="name" control={control} render={({ field }) => (
                  <TextField {...field} label="Name" fullWidth required error={!!errors.name} helperText={errors.name?.message} disabled={loading} />
                )} />
              </Grid>
              <Grid item xs={4}>
                <Controller name="shortName" control={control} render={({ field }) => (
                  <TextField {...field} label="Short Name" fullWidth error={!!errors.shortName} helperText={errors.shortName?.message} disabled={loading} />
                )} />
              </Grid>
              <Grid item xs={12}>
                <Controller name="slug" control={control} render={({ field }) => (
                  <TextField {...field} label="Slug" fullWidth error={!!errors.slug} helperText={errors.slug?.message || 'Auto-generated'} disabled={loading} />
                )} />
              </Grid>
              <Grid item xs={12}>
                <Controller name="country" control={control} render={({ field }) => (
                  <TextField {...field} label="Country" fullWidth error={!!errors.country} helperText={errors.country?.message} disabled={loading} />
                )} />
              </Grid>
              <Grid item xs={10}>
                <Controller name="logoUrl" control={control} render={({ field }) => (
                  <TextField {...field} label="Logo URL" fullWidth error={!!errors.logoUrl} helperText={errors.logoUrl?.message} disabled={loading} />
                )} />
              </Grid>
              <Grid item xs={2} sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                {logoUrl && <Avatar src={logoUrl} sx={{ width: 48, height: 48 }} />}
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

export default TeamForm
