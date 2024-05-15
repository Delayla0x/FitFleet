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
      console.error(`${message}:`, error);
    } else {
      console.log(message);
    }
  }

  private static handleError(error: any): void {
    if (error instanceof AxiosError) {
      if (error.response) {
        this.log(`Server responded with a status of ${error.response.status}`, error);
        console.error(error.response.data);
        console.error(error.response.headers);
      } else if (error.request) {
        this.log('No response received for the request', error);
        console.error(error.request);
      } else {
        this.log("Error setting up the request", error);
      }
    } else {
      this.log('Error:', error);
    }
  }

  static async fetchClassSchedules(): Promise<ClassSchedule[]> {
    try {
      const response = await axios.get(`${API_BASE_URL}/classes/schedules`);
      this.log('Successfully fetched class schedules');
      return response.data;
    } catch (error) {
      this.log('Error fetching class schedules', error);
      this.handleError(error);
      throw error;
    }
  }

  static async bookAClass(bookingRequest: BookingRequest): Promise<void> {
    try {
      await axios.post(`${API_BASE_URL}/classes/book`, bookingRequest);
      this.log('Class booking successful');
    } catch (error) {
      this.log('Error booking the class', error);
      this.handleError(error);
      throw error;
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
      this.log('Error cancelling class reservation', error);
      this.handleError(error);
      throw error;
    }
  }
}

export default ClassService;