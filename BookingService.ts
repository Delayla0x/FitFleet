import axios from 'axios';

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

  static async fetchClassSchedules(): Promise<ClassSchedule[]> {
    try {
      const response = await axios.get(`${API_BASE_URL}/classes/schedules`);
      return response.data;
    } catch (error) {
      console.error('Error fetching class schedules', error);
      throw error;
    }
  }

  static async bookAClass(bookingRequest: BookingRequest): Promise<void> {
    try {
      await axios.post(`${API_BASE_URL}/classes/book`, bookingRequest);
      console.log('Class booking successful');
    } catch (error) {
      console.error('Error booking the class', error);
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
      console.log('Class reservation cancelled successfully');
    } catch (error) {
      console.error('Error cancelling class reservation', error);
      throw error;
    }
  }
}

export default ClassService;