import React from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import styled from "styled-components";
import { useState, useEffect } from "react";
import Swal from "sweetalert2"; 
import { Navigate } from 'react-router-dom'; 

const formSchema = z.object({
  title: z.string().min(1, "Title is required"),
  overview: z.string().min(1, "Overview is required"),
  adult: z.boolean().default(false),
  originalTitle: z.string().optional(),
  originalLang: z.string().optional(),
  releaseDate: z.string().optional(),
  popularity: z.number().min(0).optional(),
});

const MovieForm = () => {
  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      title: "",
      overview: "",
      adult: false,
      originalTitle: "",
      originalLang: "",
      releaseDate: "",
      popularity: 0,
    },
  });

  const [name, setName] = useState("");
  
  useEffect(() => {
     const headers = new Headers();
     headers.append('Content-Type', 'application/json');
     const reqOptions = {
       method: 'GET',
       headers: headers,
       credentials: 'include'
     };
 
     fetch(`http://localhost:4000/Username`, reqOptions)
       .then(response => response.json())
       .then(data => setName(data.username));
 
     return () => setName(""); 
   }, []);

  const onSubmit = (values) => {
    const body = {
      title: values.title,
      backdrop_path: "",
      id: 0,
      overview: values.overview,
      adult: values.adult,
      originalTitle: values.originalTitle,
      originalLang: values.originalLang,
      releaseDate: values.releaseDate,
      poster_path: "",
      popularity: 0.0,
      vote_average: 0.0,
      votecount: 0.0,
      author: name,
    };

    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    const requestOptions = {
      method: 'POST',
      headers: headers,
      body: JSON.stringify(body),
    };

    fetch(`http://localhost:4000/addusermovies`, requestOptions)
      .then(response => response.json())
      .then(data => {
        if (data.status !==200) {
          Swal.fire({
            icon: 'error',
            title: 'Movie Insertion Failed',
            text: 'There was an issue inserting the movie.',
            confirmButtonText: 'Try Again',
          });
        } else {
          Swal.fire({
            icon: 'success',
            title: 'Movie Added Successfully',
            text: 'The movie has been successfully added!',
            confirmButtonText: 'Go to Home',
          }).then(() => {
            Navigate(`/`); 
          });
        }
      })
      .catch(err => {
        Swal.fire({
          icon: 'error',
          title: 'Something went wrong',
          text: 'There was an error with the network. Please try again later.',
          confirmButtonText: 'OK',
        });
      });
  };

  return (
    <Container>
      <FormWrapper>
        <Heading>Add New Movie</Heading>
        <SubHeading>Enter the details of the new movie below.</SubHeading>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <FormField>
            <Label>Title</Label>
            <Input
              type="text"
              placeholder="Enter movie title"
              {...form.register("title")}
            />
            <Error>{form.formState.errors.title?.message}</Error>
          </FormField>
          <FormField>
            <Label>Overview</Label>
            <TextArea
              placeholder="Enter movie overview"
              {...form.register("overview")}
            />
            <Error>{form.formState.errors.overview?.message}</Error>
          </FormField>
          <CheckboxWrapper>
            <Label>Adult Content</Label>
            <Checkbox type="checkbox" {...form.register("adult")} />
          </CheckboxWrapper>
          <FormField>
            <Label>Original Title</Label>
            <Input
              type="text"
              placeholder="Enter original title"
              {...form.register("originalTitle")}
            />
          </FormField>
          <FormField>
            <Label>Original Language</Label>
            <Input
              type="text"
              placeholder="Enter original language"
              {...form.register("originalLang")}
            />
          </FormField>
          <FormField>
            <Label>Release Date</Label>
            <Input type="date" {...form.register("releaseDate")} />
          </FormField>
          <FormField>
            <Label>Popularity</Label>
            <Input type="number" {...form.register("popularity")} />
          </FormField>
          <SubmitButton type="submit">Submit</SubmitButton>
        </Form>
      </FormWrapper>
    </Container>
  );
};


const Container = styled.div`
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #121212;
  color: #e0e0e0;
  padding: 20px;
  box-sizing: border-box;
`;

const FormWrapper = styled.div`
  width: 100%;
  max-width: 600px;
  background-color: #1e1e1e;
  padding: 50px;
  border-radius: 10px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
`;

const Heading = styled.h2`
  font-size: 26px;
  font-weight: bold;
  margin-bottom: 15px;
  color: #ffd700;
`;

const SubHeading = styled.p`
  margin-bottom: 25px;
  font-size: 16px;
  color: #a0a0a0;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

const FormField = styled.div`
  display: flex;
  flex-direction: column;
`;

const Label = styled.label`
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #cccccc;
`;

const Input = styled.input`
  width: 100%;
  padding: 12px;
  border-radius: 6px;
  border: 1px solid #444;
  background-color: #2a2a2a;
  color: #e0e0e0;
  font-size: 15px;
  transition: border-color 0.3s ease;
  &:focus {
    border-color: #ffd700;
  }
`;

const TextArea = styled.textarea`
  ${Input};
  height: 120px;
  resize: none;
`;

const CheckboxWrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Checkbox = styled.input`
  accent-color: #4caf50;
`;

const Error = styled.p`
  color: #ff4c4c;
  font-size: 12px;
  margin-top: 5px;
`;

const SubmitButton = styled.button`
  padding: 12px 20px;
  border-radius: 6px;
  background-color: #4caf50;
  color: #fff;
  font-size: 16px;
  font-weight: bold;
  border: none;
  cursor: pointer;
  text-align: center;
  &:hover {
    background-color: #45a049;
  }
`;

export default MovieForm;
