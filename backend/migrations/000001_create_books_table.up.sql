create table IF NOT EXISTS books (
    book_id bigserial PRIMARY KEY, 
    title text not null,
    author text not null,
    genres text[] not null CONSTRAINT genres_length_check check (array_length(genres, 1) between 1 and 5),
    version INTEGER not null default 1
);

insert into books (title, author, genres) values ('Emotional Intelligence', 'Daniel Goleman', array['Psychology','Science','Self-help']);
insert into books (title, author, genres) values ('Pride and Prejudice', 'Jane Austen', array['Fiction','Romance']);
insert into books (title, author, genres) values ('Lords of Finance: The Bankers Who Broke the World', 'Liaquat Ahamed', array['Non-fiction']);
insert into books (title, author, genres) values ('The Poetry Pharmacy Returns: More Prescriptions for Courage, Healing and Hope', 'William Sieghart', array['Poetry']);
insert into books (title, author, genres) values ('The Body Keeps the Score: Brain, Mind, and Body in the Healing of Trauma', 'Bessel van der Kolk', array['Non-fiction','Psychology']);
