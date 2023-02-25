import { Component, createResource, createSignal, Show, For } from 'solid-js'
import {
  Drawer,
  IconButton,
  Box,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  ListItemIcon,
  Alert
} from '@suid/material'
import EditIcon from '@suid/icons-material/Edit'
import PoolIcon from '@suid/icons-material/Pool'
import AccessTimeIcon from '@suid/icons-material/AccessTime'
import { trainingApi } from '../state/api'
import { GetDetailsCurrWeekResponse } from '../generated'
import DetailCard from '../component/DetailCard'
import { NavLink } from 'solid-app-router'

async function getTrainingsThisWeek(): Promise<GetDetailsCurrWeekResponse> {
  const result = trainingApi.getTrainingsDetailsCurrentWeek()
  const trainings = result
    .then((res) => res)
    .catch((err) => {
      throw err
    })
  return trainings
}

const Home: Component = () => {
  const [details] = createResource(getTrainingsThisWeek)
  const [open, setOpen] = createSignal(false)

  const list = () => (
    <Box
      role="presentation"
      onClick={() => setOpen(false)}
      onKeyDown={() => setOpen(false)}
    >
      <List>
        <ListItem>
          <NavLink href="/training/create">
            <ListItemButton>
              <ListItemIcon>
                <PoolIcon />
              </ListItemIcon>
              <ListItemText primary="Create Training" />
            </ListItemButton>
          </NavLink>
        </ListItem>
        <ListItem>
          <ListItemButton>
            <ListItemIcon>
              <EditIcon />
            </ListItemIcon>
            <ListItemText primary="Edit Training" />
          </ListItemButton>
        </ListItem>
        <ListItem>
          <ListItemButton>
            <ListItemIcon>
              <AccessTimeIcon />
            </ListItemIcon>
            <ListItemText primary="Create Session" />
          </ListItemButton>
        </ListItem>
        <ListItem>
          <ListItemButton>
            <ListItemIcon>
              <EditIcon />
            </ListItemIcon>
            <ListItemText primary="Edit Session" />
          </ListItemButton>
        </ListItem>
      </List>
    </Box>
  )

  return (
    <div class="mx-auto w-11/12 text-center">
      <h1 class="my-4 text-2xl font-bold">This week trainings</h1>

      <Show when={details.error}>
        <Alert severity="error">
          Couldn't load training details for this week
        </Alert>
      </Show>

      <Drawer
        anchor="right"
        open={open()}
        sx={{ zIndex: 9999 }}
        onClose={() => setOpen(false)}
      >
        {list()}
      </Drawer>

      <For each={details()?.details}>
        {(detail) => <DetailCard detail={detail} />}
      </For>

      <div
        onClick={() => setOpen(true)}
        class="absolute bottom-4 right-4 flex items-center justify-center rounded-full shadow"
      >
        <IconButton color="primary">
          <EditIcon fontSize="large" />
        </IconButton>
      </div>
    </div>
  )
}

export default Home
