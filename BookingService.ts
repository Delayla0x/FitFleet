import axios, { AxiosError } from 'axios';

interface ClassSchedule {
  id: number;
  name: string;
  date: string;
  time: string;
}

interface BookingRequest {
  classId: number;
  userId: string;
}

const API_BASE_URL: string = process.env.REACT_APP_API_BASE_URL || '';

class ClassService {
  private static log(message: string, error?: Error): void {
    if (error) {
      console.error(`${message}:`, error.message);
    } else {
      console.log(message);
    }
  }

  private static handleError(error: unknown): never {
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError;
      if (axiosError.response) {
        this.log(`Server responded with a status of ${axiosError.response.status}`, axiosError);
        console.error(axiosError.response.data);
        console.error(axiosError.response.headers);
      } else if (axiosError.request) {
        this.log('No response received for the request', axiosError);
        console.error(axiosError.request);
      } else {
        this.log("Error setting up the request", axiosError);
      }
    } else if (error instanceof Error) {
      this.log('Error:', error);
    } else {
      console.error('An unexpected error occurred:', error);
    }

    throw error;
  }

  static async fetchClassSchedules(): Promise<ClassSchedule[]> {
    try {
      const response = await axios.get(`${API_BASE_URL}/classes/schedules`);
      this.log('Successfully fetched class schedules');
      return response.data;
    } catch (error) {
      this.log('Error fetching class schedules', error as Error);
      this.handleError(error);
    }
  }

  static async bookAClass(bookingRequest: BookingRequest): Promise<void> {
    try {
      await axios.post(`${API_BASE_URL}/classes/book`, bookingRequest);
      this.log('Class booking successful');
    } catch (error) {
      this.log('Error booking the class', error as Error);
      this.handleError(error);
    }
  }

  static async cancelAClassReservation(classId: number, userId: string): Promise<void> {
    try {
      await axios.delete(`${API_BASE_URL}/classes/cancel`, {
        data: {
          classId,
          userId,
        },
      });
      this.log('Class reservation cancelled successfully');
    } catch (error) {
      this.log('Error cancelling class reservation', error as Error);
      this.handleError(error);
    }
  }
}

export default ClassService;