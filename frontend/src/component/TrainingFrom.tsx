import { Component, createSignal, Show } from 'solid-js'
import { SetStoreFunction } from 'solid-js/store'
import { Training, Day, InvalidTraining } from '../generated'
import {
  Box,
  Button,
  FormControl,
  FormHelperText,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  TextField
} from '@suid/material'
import { SelectChangeEvent } from '@suid/material/Select'

interface TraningFormProps {
  training: Training
  setTraining: SetStoreFunction<Training>
  invalidTraining: InvalidTraining
}

const [sessionManual, setSessionManual] = createSignal(true)
const [durErr, setDurErr] = createSignal(false)

const TraningForm: Component<TraningFormProps> = (props) => {
  const handleDayChange = (event: SelectChangeEvent) => {
    props.setTraining('day', event.target.value as Day)
  }

  const handleDurationChange = (value: string) => {
    const duration = parseInt(value)
    if (duration < 1) {
      setDurErr(true)
      props.setTraining('durationMin', undefined)
      return
    }

    setDurErr(false)
    props.setTraining('durationMin', duration)
  }

  return (
    <div>
      <Stack spacing={2} direction="row" justifyContent="center">
        <Button
          variant={sessionManual() ? 'contained' : 'outlined'}
          size="medium"
          onClick={() => setSessionManual(true)}
        >
          Set manually
        </Button>
        <Button
          variant={sessionManual() ? 'outlined' : 'contained'}
          size="medium"
          onClick={() => setSessionManual(false)}
        >
          Assign session
        </Button>
      </Stack>

      <Show when={sessionManual()}>
        <Box sx={{ maxWidth: '90%', mx: 'auto', my: '1em' }}>
          <FormControl
            fullWidth
            margin="normal"
            error={props.invalidTraining.day !== undefined}
          >
            <InputLabel id="day">Day</InputLabel>
            <Select
              labelId="day"
              label="Day"
              value={props.training.day ? props.training.day : ''}
              onChange={handleDayChange}
            >
              <MenuItem value={Day.Monday}>Monday</MenuItem>
              <MenuItem value={Day.Tuesday}>Tuesday</MenuItem>
              <MenuItem value={Day.Wednesday}>Wednesday</MenuItem>
              <MenuItem value={Day.Thursday}>Thursday</MenuItem>
              <MenuItem value={Day.Friday}>Friday</MenuItem>
              <MenuItem value={Day.Saturday}>Saturday</MenuItem>
              <MenuItem value={Day.Sunday}>Sunday</MenuItem>
            </Select>
            <FormHelperText>{props.invalidTraining.day}</FormHelperText>
          </FormControl>
          <FormControl
            fullWidth
            margin="normal"
            error={props.invalidTraining.durationMin !== undefined}
          >
            <TextField
              id="duration"
              label="Duration"
              type="number"
              onChange={(_, val) => handleDurationChange(val)}
              value={
                props.training.durationMin ? props.training.durationMin : null
              }
              InputProps={{
                inputProps: { min: 1 },
                endAdornment: <p>minutes</p>
              }}
              error={
                props.invalidTraining.durationMin !== undefined || durErr()
              }
            />
          </FormControl>
        </Box>
      </Show>
    </div>
  )
}

export default TraningForm
