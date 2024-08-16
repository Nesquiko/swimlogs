create table if not exists styles
(
    id                   uuid primary key,
    name                 varchar(50) not null,
    exercise_name        varchar(50),
    exercise_description text,

    unique (name, exercise_name, exercise_description)
);

insert into styles (id, name, exercise_name, exercise_description)
values ('2f9fd38b-9a56-4117-8e50-fe10feb31ae6', 'choice', null, null),

       ('d4876e51-6206-4409-8315-3fd776c3b3d0', 'medley', null, null),

       ('d4876e51-6206-4409-8315-3fd776c3b3d0', 'freestyle', null, null),
       ('d770dcc3-ba7d-4da2-9f7b-92d69fb445cd', 'freestyle', 'kick', 'freestyle.kick'),
       ('acfc0afa-50b0-4c86-bd12-55fa11eca7d6', 'freestyle', 'kick.head.up', 'freestyle.kick.head.up'),
       ('e5d01df7-194e-45aa-925e-e55786942fbb', 'freestyle', 'pull', 'freestyle.pull'),
       ('b2ab0523-eec1-408b-8ffd-78ec8d2bee18', 'freestyle', 'catch.up', 'freestyle.catch.up'),

       ('0511b862-d596-40e2-adc4-f7d003949499', 'breaststroke', null, null),
       ('76af41e7-2639-49c4-b9e4-bd4afbcf318e', 'breaststroke', 'kick', 'breaststroke.kick'),
       ('dd5ba9c5-95ee-4ad8-a5d8-2b4491cdb10d', 'breaststroke', 'pull', 'breaststroke.pull'),
       ('5be726ee-d7f4-407b-b317-c622b197d662', 'breaststroke', 'on.back', 'breaststroke.on.back'),

       ('95779412-f1b0-49a2-b775-874f090ea704', 'butterfly', null, null),
       ('d67298e2-60e0-4613-8bae-fae48f09e753', 'butterfly', 'pull', 'butterfly.pull'),
       ('cfd69122-b412-4098-b951-c6f167889eed', 'butterfly', 'dolphin.kick', 'butterfly.dolphin.kick'),

       ('89319cef-5be5-4298-a27e-483df8140715', 'backstroke', null, null),
       ('1ea55082-654d-4212-a9d4-1226bfab06eb', 'backstroke', 'kick', 'backstroke.kick'),
       ('3b937568-419b-4b98-9e17-fc0d7c5fa79d', 'backstroke', 'pull', 'backstroke.pull'),

       ('8cc9d3bc-7a0a-42df-b493-11004004e3d0', 'surface', null, null),
       ('af040fb9-e7bb-475c-b3be-492653775fcd', 'surface', 'soldier', 'surface.soldier'),
       ('7f56bc3f-3641-4e75-882a-23b9633ec162', 'surface', 'on.back', 'surface.on.back'),
       ('dc1e667e-719c-4e06-89e4-783a94a269ab', 'apnea', null, null),
       ('5a4e880e-e6d3-4390-b1f8-51cf34e13b8c', 'immersion', null, null),
       ('6528f34f-ce4b-45f5-b6bc-54f4c558ff6e', 'bifins', null, null)
on conflict do nothing;
